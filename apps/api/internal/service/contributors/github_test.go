package contributors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"contrib.rocks/libs/go/model"
	"github.com/google/go-github/v45/github"
)

func setup(t *testing.T, handler http.Handler) (*github.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})
	ghclient := github.NewClient(server.Client())
	ghclient.BaseURL, _ = url.Parse(server.URL + "/")
	return ghclient, server
}

func Test_fetchRepository(t *testing.T) {
	t.Run("should send a request to GitHub endpoint", func(t *testing.T) {
		var req *http.Request
		ghclient, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req = r
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1}`))
		}))
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}
		err := makeFetchRepositoryFn(ghclient, context.Background(), repository, nil)()
		if err != nil {
			t.Fatal(err)
		}
		if req.Method != "GET" {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/repos/foo/bar" {
			t.Fatalf("expected request to be sent to %s, got %s", "/repos/owner/repo", req.URL.String())
		}
	})

	t.Run("should return fetched repository result", func(t *testing.T) {
		ghclient, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("%s %s", r.Method, r.URL.String())
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":1}`))
		}))
		out := make(chan *github.Repository, 1)
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}
		err := makeFetchRepositoryFn(ghclient, context.Background(), repository, out)()
		if err != nil {
			t.Fatal(err)
		}
		result, ok := <-out
		if !ok {
			t.Fatal("expected result to be available")
		}
		if result.GetID() != 1 {
			t.Fatalf("expected repository id to be 1, got %d", result.GetID())
		}
	})

	t.Run("should throw not found error on 404", func(t *testing.T) {
		ghclient, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}
		err := makeFetchRepositoryFn(ghclient, context.Background(), repository, nil)()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if _, ok := err.(*RepositoryNotFoundError); !ok {
			t.Fatalf("expected RepositoryNotFoundError, got %T", err)
		}
	})
}

func Test_fetchContributors(t *testing.T) {
	t.Run("should send a request to GitHub endpoint", func(t *testing.T) {
		var req *http.Request
		ghclient, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req = r
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": 1}]`))
		}))
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}
		err := makeFetchContributorsFn(ghclient, context.Background(), repository, nil)()
		if err != nil {
			t.Fatal(err)
		}
		if req.Method != "GET" {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/repos/foo/bar/contributors" {
			t.Fatalf("expected request to be sent to %s, got %s", "/repos/owner/repo/contributors", req.URL.String())
		}
	})

	t.Run("should return fetched contributors result", func(t *testing.T) {
		ghclient, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("%s %s", r.Method, r.URL.String())
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"id": 1}]`))
		}))
		out := make(chan []*github.Contributor, 1)
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}
		err := makeFetchContributorsFn(ghclient, context.Background(), repository, out)()
		if err != nil {
			t.Fatal(err)
		}
		result, ok := <-out
		if !ok {
			t.Fatal("expected result to be available")
		}
		if len(result) != 1 {
			t.Fatalf("expected 1 contributor, got %v", result)
		}
		if result[0].GetID() != 1 {
			t.Fatalf("expected contributor id to be 1, got %d", result[0].GetID())
		}
	})

	t.Run("should throw not found error on 404", func(t *testing.T) {
		ghclient, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		repository := &model.Repository{Owner: "foo", RepoName: "bar"}
		err := makeFetchContributorsFn(ghclient, context.Background(), repository, nil)()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if _, ok := err.(*RepositoryNotFoundError); !ok {
			t.Fatalf("expected RepositoryNotFoundError, got %T", err)
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
