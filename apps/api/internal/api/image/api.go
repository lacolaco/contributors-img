package image

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	imageMaxAge = 60 * 60 * 24 * 3 // 3 days
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
		attribute.String("/app/api/image/params/referer", params.Referer),
		attribute.Int64("/app/api/image/params/max", int64(params.MaxCount)),
		attribute.Int64("/app/api/image/params/columns", int64(params.Columns)),
	)
	log = log.With(logger.Label("repository", string(params.Repository)),
		logger.Label("referer", params.Referer))
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
		sendImage(c, image, "")
		return
	}

	// get data
	data, err := api.cs.GetContributors(ctx, params.Repository.Object())
	var notfound *model.RepositoryNotFoundError
	if errors.As(err, &notfound) {
		log.Error(err.Error())
		// Cache the 404 for 1 hour to prevent GitHub API exhaustion
		c.Header("Cache-Control", "public, max-age=3600, s-maxage=3600")
		c.String(http.StatusNotFound, notfound.Error())
		return
	} else if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	var dataETag string
	if jsonData, err := json.Marshal(data); err == nil {
		hash := md5.Sum(jsonData)
		dataETag = fmt.Sprintf(`"%x"`, hash)

		if strings.Contains(c.GetHeader("If-None-Match"), fmt.Sprintf("%x", hash)) {
			c.Status(http.StatusNotModified)
			c.Header("Cache-Control", "public, max-age=43200, s-maxage=43200, stale-while-revalidate=86400")
			return
		}
	}

	// render image
	image, err = api.is.RenderImage(ctx, data, rendererOptions, params.IncludeAnonymous)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	api.us.CollectUsage(ctx, data, params.Via)
	sendImage(c, image, dataETag)
}

func sendImage(c *gin.Context, image model.FileHandle, etag string) {
	if etag == "" {
		etag = image.ETag()
	}

	// Handle Weak ETags (W/"...") and comma-separated lists
	if strings.Contains(c.GetHeader("If-None-Match"), strings.Trim(etag, `W/"`)) {
		c.Status(http.StatusNotModified)
		c.Header("Cache-Control", "public, max-age=43200, s-maxage=43200, stale-while-revalidate=86400")
		return
	}
	r := image.Reader()
	defer r.Close()
	c.DataFromReader(http.StatusOK, image.Size(), image.ContentType(), r, map[string]string{
		"Cache-Control": "public, max-age=43200, s-maxage=43200, stale-while-revalidate=86400",
		"ETag":          etag,
	})
}
