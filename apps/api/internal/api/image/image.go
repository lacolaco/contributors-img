package image

import (
	"fmt"
	"net/http"
	"strings"

	"contrib.rocks/apps/api/internal/service"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
	"github.com/gin-gonic/gin"
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
	file, err := c.is.GetImage(ctx, data, &renderer.RendererOptions{
		MaxCount: params.MaxCount,
		Columns:  params.Columns,
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
