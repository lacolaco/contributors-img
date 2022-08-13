package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"contrib.rocks/apps/api/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestMiddleware(t *testing.T) {
	var logger *zap.Logger
	r := gin.New()
	r.Use(MiddlewareZap(config.NewTestConfig()))
	r.GET("/ping", func(c *gin.Context) {
		logger = LoggerFromContext(c.Request.Context())
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if logger == nil {
		t.Fatalf("Expected logger to be set in context")
	}
}
