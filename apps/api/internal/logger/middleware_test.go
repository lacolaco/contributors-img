package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware(t *testing.T) {
	logger := newStdLogger("test")

	var got Logger
	r := gin.New()
	r.Use(Middleware(logger))
	r.GET("/ping", func(c *gin.Context) {
		got = FromContext(c.Request.Context())
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if got == nil {
		t.Fatalf("Expected logger to be set in context")
	}
	if got != logger {
		t.Fatalf("Expected logger to be set in context")
	}
}
