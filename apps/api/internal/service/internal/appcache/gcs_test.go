package appcache

import (
	"errors"
	"fmt"
	"testing"

	"cloud.google.com/go/storage"
)

// A cache miss must stay a miss. storage returned ErrObjectNotExist bare until
// v1.51.0 and wraps it since, so an equality check silently turns every miss
// into a request failure — which no other test here can observe, because the
// in-memory cache cannot return an error at all.
func Test_isNotExistErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"bare sentinel", storage.ErrObjectNotExist, true},
		{"wrapped sentinel", fmt.Errorf("%w: %w", storage.ErrObjectNotExist, errors.New("googleapi: got HTTP response code 404")), true},
		{"doubly wrapped sentinel", fmt.Errorf("reading object: %w", fmt.Errorf("%w: %w", storage.ErrObjectNotExist, errors.New("404"))), true},
		{"unrelated error", errors.New("connection reset by peer"), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNotExistErr(tt.err); got != tt.want {
				t.Fatalf("isNotExistErr(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
