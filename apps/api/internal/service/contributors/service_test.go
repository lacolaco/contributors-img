package contributors

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/internal/logger"
	"github.com/google/go-github/v69/github"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type stubProvider struct{ client *github.Client }

func (p *stubProvider) Get() *github.Client { return p.client }

// mockCache reports a miss on read and a failure on write.
type mockCache struct{ saveJSONErr error }

func (m *mockCache) Get(context.Context, string) (model.FileHandle, error) { return nil, nil }
func (m *mockCache) GetJSON(context.Context, string, any) error            { return nil }
func (m *mockCache) Save(context.Context, string, []byte, string) error    { return nil }
func (m *mockCache) SaveJSON(context.Context, string, any) error           { return m.saveJSONErr }

// A cache write that fails must not fail the request. The data in hand was paid
// for with GitHub API quota — the scarcest resource this service has — so
// discarding it because GCS refused a write means the next request buys it
// again. Worse, a persistent write failure would then burn the hourly limit and
// take the service down for a reason unrelated to the actual fault.
func TestService_GetContributors_cacheSaveFailureIsNotFatal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/test-owner/test-repo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"stargazers_count":42,"owner":{"login":"test-owner"},"name":"test-repo"}`))
	})
	mux.HandleFunc("/repos/test-owner/test-repo/contributors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"login":"user1","avatar_url":"https://avatar1.example"}]`))
	})
	ghclient, _ := setup(t, mux)

	service := New(&stubProvider{client: ghclient}, &mockCache{
		saveJSONErr: errors.New("googleapi: Error 403: quota exceeded"),
	})

	core, logs := observer.New(zapcore.DebugLevel)
	ctx := logger.ContextWithLogger(context.Background(), zap.New(core))

	data, err := service.GetContributors(ctx, &model.Repository{Owner: "test-owner", RepoName: "test-repo"})
	if err != nil {
		t.Fatalf("err = %v, want nil: a failed cache write must not fail the request", err)
	}
	if data == nil || len(data.Contributors) != 1 {
		t.Fatalf("data = %+v, want the contributors fetched from GitHub", data)
	}

	// The write still has to be reported, or the failure is invisible again.
	errorLogs := logs.FilterLevelExact(zapcore.ErrorLevel).All()
	if len(errorLogs) != 1 {
		t.Fatalf("error logs = %d, want 1: the failure must still be recorded", len(errorLogs))
	}
	if !strings.Contains(errorLogs[0].Message, "contributors-json-cache-save-failure") {
		t.Fatalf("log = %q, want it to name the cache-save-failure group", errorLogs[0].Message)
	}
}
