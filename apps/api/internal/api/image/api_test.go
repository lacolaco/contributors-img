package image

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
	"contrib.rocks/apps/api/internal/apierror"
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

// stubImageService always reports a cache miss so that Get reaches GetContributors,
// matching what internal/service/image returns when the cache has no entry.
type stubImageService struct{}

func (s *stubImageService) GetImage(context.Context, *model.Repository, *renderer.RendererOptions, bool) (model.FileHandle, error) {
	return nil, nil
}

// RenderImage is never reached by these cases; fail loudly rather than returning
// a nil handle that would panic further down in sendImage.
func (s *stubImageService) RenderImage(context.Context, *model.RepositoryContributors, *renderer.RendererOptions, bool) (model.FileHandle, error) {
	return nil, errors.New("RenderImage should not be reached")
}

type stubUsageService struct{}

func (s *stubUsageService) CollectUsage(context.Context, *model.RepositoryContributors, string) error {
	return nil
}

// newTestRouter wires the handler the way internal/server.go does, minus the
// middleware that needs external services.
func newTestRouter(contributorsErr error) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.Middleware(config.NewTestConfig()))
	r.Use(apierror.Middleware())
	api := New(&stubContributorsService{err: contributorsErr}, &stubImageService{}, &stubUsageService{})
	r.GET("/image", api.Get)
	return r
}

// A missing repository must surface as 404, not 500 — and everything else must
// stay 500, since the match is by errors.As rather than by exact type.
//
// The not-found error does not arrive bare. HandleRepositoryNotFoundError returns
// retry.Unrecoverable, and retry.Do records unpackUnrecoverable(err) — so the
// Unrecoverable wrapper is stripped, but the surrounding retry.Error is not,
// because GetRetryOptions does not set LastErrorOnly. The handler therefore sees
// retry.Error{notFound}, which a plain type assertion misses.
//
// The wrapper is built by retry.Do rather than by hand so that this fixture
// cannot drift from what production actually produces.
func Test_Get_RepositoryNotFound(t *testing.T) {
	repo := &model.Repository{Owner: "lacolaco", RepoName: "does-not-exist"}
	notFound := &model.RepositoryNotFoundError{Repository: repo}

	wrapped := retry.Do(func() error { return retry.Unrecoverable(notFound) })
	if wrapped == nil {
		t.Fatal("retry.Do returned nil, expected the error it was given")
	}

	tests := []struct {
		name string
		err  error
		want int
		// body must appear in the response, so that a status reached for the wrong
		// reason — a panic recovered as 500, say — does not pass for the right one.
		body string
	}{
		{"bare", notFound, http.StatusNotFound, "lacolaco/does-not-exist"},
		{"wrapped by retry.Do", wrapped, http.StatusNotFound, "lacolaco/does-not-exist"},
		{"unrelated error stays 500", errors.New("github is unreachable"), http.StatusInternalServerError, "github is unreachable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/image?repo=lacolaco/does-not-exist", nil)
			newTestRouter(tt.err).ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d (body: %s)", w.Code, tt.want, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tt.body) {
				t.Fatalf("body = %q, want it to contain %q", w.Body.String(), tt.body)
			}
		})
	}
}
