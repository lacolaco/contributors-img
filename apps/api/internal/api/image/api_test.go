package image

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/logger"
	"github.com/avast/retry-go/v4"
	"github.com/gin-gonic/gin"
)

// stubContributorsService returns a fixed error from GetContributors.
type stubContributorsService struct{ err error }

func (s *stubContributorsService) GetContributors(context.Context, *model.Repository) (*model.RepositoryContributors, error) {
	return nil, s.err
}

// stubImageService always reports a cache miss so that Get reaches GetContributors.
type stubImageService struct{}

func (s *stubImageService) GetImage(context.Context, *model.Repository, *renderer.RendererOptions, bool) (model.FileHandle, error) {
	return nil, nil
}

func (s *stubImageService) RenderImage(context.Context, *model.RepositoryContributors, *renderer.RendererOptions, bool) (model.FileHandle, error) {
	return nil, nil
}

type stubUsageService struct{}

func (s *stubUsageService) CollectUsage(context.Context, *model.RepositoryContributors, string) error {
	return nil
}

// newTestRouter wires the handler the way internal/server.go does, minus the
// middleware that needs external services. errorHandler is reproduced here
// because it is unexported in the parent package; keep it in step with it.
func newTestRouter(contributorsErr error) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(logger.Middleware(config.NewTestConfig()))
	r.Use(func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
		}
	})
	api := New(&stubContributorsService{err: contributorsErr}, &stubImageService{}, &stubUsageService{})
	r.GET("/image", api.Get)
	return r
}

// A missing repository must surface as 404, not 500. The error does not arrive
// bare: HandleRepositoryNotFoundError returns retry.Unrecoverable, and
// retry.Do wraps that in a retry.Error unless LastErrorOnly is set — which
// GetRetryOptions does not set. A plain type assertion misses it and the
// request falls through to the 500 handler.
func Test_Get_RepositoryNotFound(t *testing.T) {
	repo := &model.Repository{Owner: "lacolaco", RepoName: "does-not-exist"}
	notFound := &model.RepositoryNotFoundError{Repository: repo}

	tests := []struct {
		name string
		err  error
	}{
		{"bare", notFound},
		{"wrapped by retry.Do", retry.Error{retry.Unrecoverable(notFound)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/image?repo=lacolaco/does-not-exist", nil)
			newTestRouter(tt.err).ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d (body: %s)", w.Code, http.StatusNotFound, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), "lacolaco/does-not-exist") {
				t.Fatalf("body = %q, want it to name the repository", w.Body.String())
			}
		})
	}
}
