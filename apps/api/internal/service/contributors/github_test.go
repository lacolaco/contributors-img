package contributors

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/internal/github/api"
	"github.com/avast/retry-go/v4"
	"github.com/google/go-github/v69/github"
)

func unwrapError(err error) error {
	if err == nil {
		return nil
	}

	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		return unwrapped
	}

	return err
}

func setup(t *testing.T, handler http.Handler) (*github.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})
	ghclient := github.NewClient(server.Client())
	ghclient.BaseURL, _ = url.Parse(server.URL + "/")
	return ghclient, server
}

func Test_fetchRepositoryContributors(t *testing.T) {
	t.Run("should fetch repository and contributors data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/repos/foo/bar" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"id": 1, "name": "bar", "owner": {"login": "foo"}, "stargazers_count": 100}`))
				return
			}

			if r.URL.Path == "/repos/foo/bar/contributors" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`[{"id": 1, "login": "user1"}, {"id": 2, "login": "user2"}]`))
				return
			}

			t.Fatalf("unexpected request to %s", r.URL.Path)
		})

		ghclient, _ := setup(t, handler)
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}

		result, err := fetchRepositoryContributors(ghclient, context.Background(), repository)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify repository data
		if result.Owner != "foo" {
			t.Errorf("expected owner to be foo, got %s", result.Owner)
		}
		if result.RepoName != "bar" {
			t.Errorf("expected repo name to be bar, got %s", result.RepoName)
		}
		if result.StargazersCount != 100 {
			t.Errorf("expected stargazers to be 100, got %d", result.StargazersCount)
		}

		// Verify contributors data
		if len(result.Contributors) != 2 {
			t.Fatalf("expected 2 contributors, got %d", len(result.Contributors))
		}
		if result.Contributors[0].ID != 1 {
			t.Errorf("expected contributor ID to be 1, got %d", result.Contributors[0].ID)
		}
		if result.Contributors[1].ID != 2 {
			t.Errorf("expected contributor ID to be 2, got %d", result.Contributors[1].ID)
		}
	})

	t.Run("should handle not found error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		ghclient, _ := setup(t, handler)
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}

		_, err := fetchRepositoryContributors(ghclient, context.Background(), repository)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		err = unwrapError(err)
		var repoNotFoundErr *model.RepositoryNotFoundError
		if !errors.As(err, &repoNotFoundErr) {
			t.Fatalf("expected RepositoryNotFoundError, got %T: %v", err, err)
		}
	})
}

func Test_buildRepositoryContributors(t *testing.T) {
	t.Run(".Owner should equal to repo.owner.login", func(t *testing.T) {
		rawRepo := &github.Repository{
			Owner: &github.User{Login: github.String("foo")},
		}
		rawContribs := []*github.Contributor{}
		got := buildRepositoryContributors(rawRepo, rawContribs)
		if got.Owner != "foo" {
			t.Fatalf("expected owner to be foo, got %s", got.Owner)
		}
	})
	t.Run(".RepoName should equal to repo.name", func(t *testing.T) {
		rawRepo := &github.Repository{Name: github.String("bar")}
		rawContribs := []*github.Contributor{}
		got := buildRepositoryContributors(rawRepo, rawContribs)
		if got.RepoName != "bar" {
			t.Fatalf("expected repo name to be bar, got %s", got.RepoName)
		}
	})
	t.Run(".Stargazors should equal to repo.stargazers_count", func(t *testing.T) {
		rawRepo := &github.Repository{StargazersCount: github.Int(1)}
		rawContribs := []*github.Contributor{}
		got := buildRepositoryContributors(rawRepo, rawContribs)
		if got.StargazersCount != 1 {
			t.Fatalf("expected stargazers to be 1, got %d", got.StargazersCount)
		}
	})
	t.Run(".Contributors should equal to contributors", func(t *testing.T) {
		rawRepo := &github.Repository{}
		rawContribs := []*github.Contributor{
			{ID: github.Int64(1)},
			{ID: github.Int64(2)},
		}
		got := buildRepositoryContributors(rawRepo, rawContribs)
		if len(got.Contributors) != 2 {
			t.Fatalf("expected 2 contributors, got %d", len(got.Contributors))
		}
		if got.Contributors[0].ID != 1 {
			t.Fatalf("expected contributor id to be 1, got %d", got.Contributors[0].ID)
		}
		if got.Contributors[1].ID != 2 {
			t.Fatalf("expected contributor id to be 2, got %d", got.Contributors[1].ID)
		}
	})
	t.Run(".Contributors should ignore bot users", func(t *testing.T) {
		rawRepo := &github.Repository{}
		rawContribs := []*github.Contributor{
			{ID: github.Int64(1)},
			{ID: github.Int64(2)},
			{ID: github.Int64(3), Type: github.String("Bot")},
		}
		got := buildRepositoryContributors(rawRepo, rawContribs)
		if len(got.Contributors) != 2 {
			t.Fatalf("expected 2 contributors, got %d", len(got.Contributors))
		}
		if got.Contributors[0].ID != 1 {
			t.Fatalf("expected contributor id to be 1, got %d", got.Contributors[0].ID)
		}
		if got.Contributors[1].ID != 2 {
			t.Fatalf("expected contributor id to be 2, got %d", got.Contributors[1].ID)
		}
	})
	t.Run(".Contributors should normarize anonymous users", func(t *testing.T) {
		rawRepo := &github.Repository{}
		rawContribs := []*github.Contributor{
			{Type: github.String("Anonymous"), Name: github.String("foo"), Email: github.String("foo@example.com")},
		}
		got := buildRepositoryContributors(rawRepo, rawContribs)
		if len(got.Contributors) != 1 {
			t.Fatalf("expected 1 contributor, got %d", len(got.Contributors))
		}
		if got.Contributors[0].Login != "foo" {
			t.Fatalf("unexpected contributor login, got %s", got.Contributors[0].Login)
		}
		if !strings.HasPrefix(got.Contributors[0].AvatarURL, "https://www.gravatar.com/avatar/") {
			t.Fatalf("unexpected contributor avatar url, got %s", got.Contributors[0].AvatarURL)
		}
	})
}

// テスト用にローカル版のエラー型関数を定義
func getGitHubRetryOptions(ctx context.Context, repo *model.Repository) []retry.Option {
	return api.GetRetryOptions(ctx)
}

func Test_getGitHubRetryOptions(t *testing.T) {
	repo := &model.Repository{Owner: "test", RepoName: "test-repo"}
	ctx := context.Background()
	options := getGitHubRetryOptions(ctx, repo)

	// 具体的な数値のチェックを避け、単にオプションが返されることを確認
	if len(options) == 0 {
		t.Fatalf("expected retry options, got empty slice")
	}
}

func Test_getGitHubErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		resp     *github.Response
		expected api.ErrorType
	}{
		{
			name:     "nil error and nil response",
			err:      nil,
			resp:     nil,
			expected: api.ErrorTypeUnknown,
		},
		{
			name:     "nil error with NotFound response",
			err:      nil,
			resp:     &github.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			expected: api.ErrorTypeNotFound,
		},
		{
			name:     "timeout error",
			err:      &net.DNSError{IsTimeout: true},
			resp:     nil,
			expected: api.ErrorTypeTimeout,
		},
		{
			name:     "server error (500)",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 500}},
			resp:     nil,
			expected: api.ErrorTypeServer,
		},
		{
			name:     "server error (503)",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 503}},
			resp:     nil,
			expected: api.ErrorTypeServer,
		},
		{
			name:     "not found error via ErrorResponse",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 404}},
			resp:     nil,
			expected: api.ErrorTypeNotFound,
		},
		{
			name:     "rate limit error",
			err:      &github.RateLimitError{},
			resp:     nil,
			expected: api.ErrorTypeRateLimit,
		},
		{
			name:     "abuse rate limit error",
			err:      &github.AbuseRateLimitError{},
			resp:     nil,
			expected: api.ErrorTypeAbuseRateLimit,
		},
		{
			name:     "connection refused error via OpError",
			err:      &net.OpError{Op: "dial", Err: errors.New("connection refused")},
			resp:     nil,
			expected: api.ErrorTypeConnectionRefused,
		},
		{
			name:     "connection refused error via url.Error",
			err:      &url.Error{Op: "get", URL: "http://example.com", Err: errors.New("connection refused")},
			resp:     nil,
			expected: api.ErrorTypeConnectionRefused,
		},
		{
			name:     "unauthorized error (401)",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 401}},
			resp:     nil,
			expected: api.ErrorTypeUnauthorized,
		},
		{
			name:     "forbidden error (403)",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 403}},
			resp:     nil,
			expected: api.ErrorTypeForbidden,
		},
		{
			name:     "client error (422)",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 422}},
			resp:     nil,
			expected: api.ErrorTypeClientError,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			resp:     nil,
			expected: api.ErrorTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := api.GetErrorType(tt.err, tt.resp)
			if got != tt.expected {
				t.Errorf("GetErrorType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_isRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "timeout error",
			err:      &net.DNSError{IsTimeout: true},
			expected: true,
		},
		{
			name:     "server error",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 500}},
			expected: true,
		},
		{
			name:     "connection refused error",
			err:      &net.OpError{Op: "dial", Err: errors.New("connection refused")},
			expected: true,
		},
		{
			name:     "rate limit error",
			err:      &github.RateLimitError{},
			expected: true,
		},
		{
			name:     "abuse rate limit error",
			err:      &github.AbuseRateLimitError{},
			expected: true,
		},
		{
			name:     "not found error",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 404}},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := api.IsRetryableError(tt.err)
			if got != tt.expected {
				t.Errorf("IsRetryableError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_isNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		resp     *github.Response
		expected bool
	}{
		{
			name:     "nil error and nil response",
			err:      nil,
			resp:     nil,
			expected: false,
		},
		{
			name:     "nil error with NotFound response",
			err:      nil,
			resp:     &github.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			expected: true,
		},
		{
			name:     "not found error via ErrorResponse",
			err:      &github.ErrorResponse{Response: &http.Response{StatusCode: 404}},
			resp:     nil,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			resp:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := api.IsNotFoundError(tt.err, tt.resp)
			if got != tt.expected {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.expected)
			}
		})
	}
}
