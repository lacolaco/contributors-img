package image

import (
	"context"
	"fmt"
	"net/http"

	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service/contributors"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/go/model"
	"contrib.rocks/libs/go/renderer"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
// imageMaxAge = 60 * 60 * 24 * 3 // 3 days
)

type ImageService interface {
	GetImage(ctx context.Context, repo *model.Repository, options *renderer.RendererOptions, includeAnonymous bool) (model.FileHandle, error)
	RenderImage(ctx context.Context, data *model.RepositoryContributors, options *renderer.RendererOptions, includeAnonymous bool) (model.FileHandle, error)
}

type ContributorsService interface {
	GetContributors(ctx context.Context, repo *model.Repository) (*model.RepositoryContributors, error)
}

type UsageService interface {
	CollectUsage(c context.Context, r *model.RepositoryContributors, via string) error
}

type API struct {
	cs ContributorsService
	is ImageService
	us UsageService
}

func New(cs ContributorsService, is ImageService, us UsageService) *API {
	return &API{cs, is, us}
}

func (api *API) Get(c *gin.Context) {
	ctx, span := tracing.Tracer().Start(c.Request.Context(), "api.image.Get")
	defer span.End()
	log := logger.LoggerFromContext(ctx)
	var params GetImageParams
	if err := params.bind(c); err != nil {
		log.Error(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	span.SetAttributes(
		attribute.String("/app/api/image/params/repository", string(params.Repository)),
		attribute.String("/app/api/image/params/via", params.Via),
		attribute.Int64("/app/api/image/params/max", int64(params.MaxCount)),
		attribute.Int64("/app/api/image/params/columns", int64(params.Columns)),
	)
	log = log.With(logger.Label("repository", string(params.Repository)))
	ctx = logger.ContextWithLogger(ctx, log)

	log.Info(fmt.Sprintf("[api.image.Get] start: %s", params.Repository), zap.Object("params", params))
	defer log.Info(fmt.Sprintf("[api.image.Get] end: %s", params.Repository))

	var image model.FileHandle
	rendererOptions := &renderer.RendererOptions{
		MaxCount: params.MaxCount,
		Columns:  params.Columns,
	}

	// get image
	image, err := api.is.GetImage(ctx, params.Repository.Object(), rendererOptions, params.IncludeAnonymous)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	if image != nil {
		sendImage(c, image)
		return
	}

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

	// render image
	image, err = api.is.RenderImage(ctx, data, rendererOptions, params.IncludeAnonymous)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	api.us.CollectUsage(ctx, data, params.Via)
	sendImage(c, image)
}

func sendImage(c *gin.Context, image model.FileHandle) {
	c.Status(http.StatusFound)
	c.Header("Location", image.DownloadURL())
}
