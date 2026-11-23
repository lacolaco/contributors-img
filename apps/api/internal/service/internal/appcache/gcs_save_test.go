package appcache

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// newBucketAgainst points a storage client at a local test server, so that the
// upload path can be exercised without credentials or a real bucket.
func newBucketAgainst(t *testing.T, h http.Handler) *storage.BucketHandle {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	client, err := storage.NewClient(context.Background(),
		option.WithEndpoint(srv.URL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })
	return client.Bucket("test-bucket")
}

// A failed upload must reach the caller. storage.Writer.Write only fills an
// internal pipe and reports nothing about the transfer; the result is delivered
// by Close. Discarding Close's error makes every write failure — credentials,
// quota, network, checksum — look like a successful cache save, so the entry is
// missed on every subsequent request with nothing logged and nothing returned.
func Test_saveFile_uploadFailure(t *testing.T) {
	bucket := newBucketAgainst(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"code":403,"message":"quota exceeded"}}`, http.StatusForbidden)
	}))

	err := saveFile(bucket, context.Background(), "image/foo--bar.svg", []byte("<svg/>"), "image/svg+xml")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Fatalf("err = %v, want it to mention the 403 from the server", err)
	}
}

// A nil bucket means the cache is disabled, not that saving failed.
func Test_saveFile_nilBucket(t *testing.T) {
	if err := saveFile(nil, context.Background(), "image/foo--bar.svg", []byte("<svg/>"), "image/svg+xml"); err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}
