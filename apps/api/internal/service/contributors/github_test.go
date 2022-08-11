package contributors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"contrib.rocks/libs/goutils/model"
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
