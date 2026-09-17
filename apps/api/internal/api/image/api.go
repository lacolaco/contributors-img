package image

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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
	// GetContributors returns the repository's contributors.
	//
	// When the repository does not exist, the returned error must match
	// *model.RepositoryNotFoundError under errors.As. When the GitHub API rate
	// limit is exhausted, it must match *model.RateLimitedError instead.
	// Implementations may wrap either — the production one returns them inside
	// a retry.Error — but must not replace them with an opaque error: Get maps
	// the first case to 404, the second to 503, and everything else to 500.
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
		sendImage(c, image)
		return
	}

	// get data
	data, err := api.cs.GetContributors(ctx, params.Repository.Object())
	// retry-go wraps this in a retry.Error unless LastErrorOnly is set, so match through the wrapper.
	var notfound *model.RepositoryNotFoundError
	var rateLimited *model.RateLimitedError
	if errors.As(err, &notfound) {
		log.Error(err.Error())
		c.String(http.StatusNotFound, notfound.Error())
		return
	} else if errors.As(err, &rateLimited) {
		log.Warn(err.Error())
		sendRateLimited(c, rateLimited)
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

// sendRateLimited answers a GitHub API rate limit with 503 rather than 500,
// and caches the response for the same duration it asks the client to wait.
// This endpoint is fetched by GitHub's camo proxy on behalf of every README
// that embeds the image, so a cacheable error gives any cache in front of us
// the chance to stop re-requesting the same repository while the rate limit
// is still exhausted.
func sendRateLimited(c *gin.Context, err *model.RateLimitedError) {
	retryAfterSeconds := int(err.RetryAfter.Seconds())
	if retryAfterSeconds <= 0 {
		retryAfterSeconds = 60
	}
	c.Header("retry-after", fmt.Sprintf("%d", retryAfterSeconds))
	c.Header("cache-control", fmt.Sprintf("public, max-age=%d", retryAfterSeconds))
	c.String(http.StatusServiceUnavailable, err.Error())
}

func sendImage(c *gin.Context, image model.FileHandle) {
	if c.GetHeader("If-None-Match") == image.ETag() {
		c.Status(http.StatusNotModified)
		c.Header("cache-control", fmt.Sprintf("public, max-age=%d", imageMaxAge))
		return
	}
	r := image.Reader()
	defer r.Close()
	c.DataFromReader(http.StatusOK, image.Size(), image.ContentType(), r, map[string]string{
		"cache-control": fmt.Sprintf(`public, max-age=%d`, imageMaxAge),
		"etag":          image.ETag(),
	})
}
