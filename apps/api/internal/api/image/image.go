package image

import (
	"fmt"
	"net/http"
	"strings"

	"contrib.rocks/apps/api/internal/service/contributors"
	"contrib.rocks/apps/api/internal/service/image"
	"contrib.rocks/apps/api/internal/service/usage"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
	"github.com/gin-gonic/gin"
)

type API struct {
	cs *contributors.Service
	is *image.Service
	us *usage.Service
}

func New(cs *contributors.Service, is *image.Service, us *usage.Service) *API {
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

func (c *API) Get(ctx *gin.Context) {
	var params getImageParams
	if err := params.bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("params=%+v\n", params)

	// get data
	data, err := c.cs.GetContributors(ctx, params.Repository.Object())
	if err != nil {
		ctx.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	fmt.Printf("data=%+v\n", data)
	// get image
	file, err := c.is.GetImage(ctx, &image.GetImageParams{
		Repository: params.Repository.Object(),
		RendererOptions: &renderer.RendererOptions{
			MaxCount: params.MaxCount,
			Columns:  params.Columns,
		},
		Data: data,
	})
	if err != nil {
		ctx.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	r := file.Reader()
	defer r.Close()
	ctx.DataFromReader(http.StatusOK, file.Size(), file.ContentType(), r, map[string]string{
		"cache-control": fmt.Sprintf(`public, max-age=%d`, 60*60*6),
	})
	// collect usage stats
	if err := c.us.CollectUsage(ctx, data, params.Via); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypePrivate)
	}
}
