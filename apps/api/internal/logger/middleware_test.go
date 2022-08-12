package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware(t *testing.T) {
	var logger Logger
	var factory LoggerFactory
	r := gin.New()
	r.Use(Middleware(newStdLoggerFactory(), "test"))
	r.GET("/ping", func(c *gin.Context) {
		logger = LoggerFromContext(c.Request.Context())
		factory = LoggerFactoryFromContext(c.Request.Context())
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if logger == nil {
		t.Fatalf("Expected logger to be set in context")
	}
	if factory == nil {
		t.Fatalf("Expected logger factory to be set in context")
	}
}
