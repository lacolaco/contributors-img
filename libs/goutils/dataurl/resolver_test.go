package dataurl_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"contrib.rocks/libs/goutils/dataurl"
	"github.com/bradleyjkemp/cupaloy"
)

func TestDataURLResolver(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		fmt.Fprint(w, "<svg></svg>")
	}))
	defer ts.Close()

	ret, err := dataurl.ResolveImageDataURL(ts.URL, 64)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cupaloy.SnapshotT(t, ret)
}
