package image

import (
	"fmt"
	"net/http"
	"strings"

	"contrib.rocks/apps/api/internal/service"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type API struct {
	cs service.ContributorsService
	is service.ImageService
	us service.UsageService
}

func New(cs service.ContributorsService, is service.ImageService, us service.UsageService) *API {
	return &API{cs, is, us}
}

type getImageParams struct {
	Repository model.RepositoryString `form:"repo" binding:"required"`
	MaxCount   int                    `form:"max"`
	Columns    int                    `form:"columns"`
	Preview    bool                   `form:"preview"`
	Via        string
}

func (p *getImageParams) bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindQuery(p); err != nil {
		return err
	}
	// validate repository name format
	if err := model.ValidateRepositoryName(string(p.Repository)); err != nil {
		return err
	}
	p.Via = "unknown"
	if strings.Contains(ctx.Request.UserAgent(), "github") {
		p.Via = "github"
	} else if strings.HasSuffix(ctx.Request.Host, "contrib.rocks") {
		p.Via = "preview"
	}
	return nil
}

func (api *API) Get(c *gin.Context) {
	ctx, span := tracing.DefaultTracer.Start(c.Request.Context(), "api.image.Get")
	defer span.End()

	var params getImageParams
	if err := params.bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("params=%+v\n", params)
	span.SetAttributes(
		attribute.String("api.image.params.repository", string(params.Repository)),
		attribute.String("api.image.params.via", params.Via),
		attribute.Int64("api.image.params.max", int64(params.MaxCount)),
		attribute.Int64("api.image.params.columns", int64(params.Columns)),
	)

	// get data
	data, err := api.cs.GetContributors(ctx, params.Repository.Object())
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	fmt.Printf("data=%+v\n", data)
	// get image
	file, err := api.is.GetImage(ctx, data, &renderer.RendererOptions{
		MaxCount: params.MaxCount,
		Columns:  params.Columns,
	})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	r := file.Reader()
	defer r.Close()
	c.DataFromReader(http.StatusOK, file.Size(), file.ContentType(), r, map[string]string{
		"cache-control": fmt.Sprintf(`public, max-age=%d`, 60*60*6),
		"etag":          file.Etag(),
	})
	// collect usage stats
	if err := api.us.CollectUsage(ctx, data, params.Via); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
	}
}
