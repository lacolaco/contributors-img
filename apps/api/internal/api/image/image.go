package image

import (
	"fmt"
	"net/http"
	"strings"

	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service"
	"contrib.rocks/apps/api/internal/service/contributors"
	"contrib.rocks/apps/api/internal/service/image"
	"contrib.rocks/apps/api/internal/service/usage"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

const (
	imageMaxAge = 60 * 60 * 6 // 6 hours
)

type API struct {
	cs contributors.Service
	is image.Service
	us usage.Service
}

func New(sp *service.ServicePack) *API {
	return &API{sp.ContributorsService, sp.ImageService, sp.UsageService}
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
	ctx, span := tracing.Tracer().Start(c.Request.Context(), "api.image.Get")
	defer span.End()

	log := logger.LoggerFromContext(ctx)

	var params getImageParams
	if err := params.bind(c); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log = log.With(logger.Label("repository", string(params.Repository)))
	ctx = logger.ContextWithLogger(ctx, log)

	log.Sugar().Debug(params)
	span.SetAttributes(
		attribute.String("api.image.params.repository", string(params.Repository)),
		attribute.String("api.image.params.via", params.Via),
		attribute.Int64("api.image.params.max", int64(params.MaxCount)),
		attribute.Int64("api.image.params.columns", int64(params.Columns)),
	)

	// get data
	data, err := api.cs.GetContributors(ctx, params.Repository.Object())
	if notfound, ok := err.(*contributors.RepositoryNotFoundError); ok {
		log.Error(err.Error())
		c.String(http.StatusNotFound, notfound.Error())
		return
	} else if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
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
	// send image
	if c.GetHeader("If-None-Match") == file.ETag() {
		c.Status(http.StatusNotModified)
		c.Header("cache-control", fmt.Sprintf("public, max-age=%d", imageMaxAge))
		return
	}
	c.DataFromReader(http.StatusOK, file.Size(), file.ContentType(), r, map[string]string{
		"cache-control": fmt.Sprintf(`public, max-age=%d`, imageMaxAge),
		"etag":          file.ETag(),
	})
	// collect usage stats
	if err := api.us.CollectUsage(ctx, data, params.Via); err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
	}
}
