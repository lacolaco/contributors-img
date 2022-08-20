package compress

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCompress(t *testing.T) {
	t.Run("No compression with no accepted encoding", func(t *testing.T) {
		r := gin.New()
		r.Use(Compress())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "test")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept-Encoding", "")
		r.ServeHTTP(w, req)

		if w.Header().Get("Content-Encoding") != "" {
			t.Fatalf("Expected no content encoding")
		}
	})
	t.Run("Support gzip compression", func(t *testing.T) {
		r := gin.New()
		r.Use(Compress())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "test")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		r.ServeHTTP(w, req)

		if w.Header().Get("Content-Encoding") != "gzip" {
			t.Fatalf("Unsupported content encoding %s", w.Header().Get("Content-Encoding"))
		}
	})
	t.Run("Support brotli compression", func(t *testing.T) {
		r := gin.New()
		r.Use(Compress())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "test")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept-Encoding", "br")
		r.ServeHTTP(w, req)

		if w.Header().Get("Content-Encoding") != "br" {
			t.Fatalf("Unsupported content encoding %s", w.Header().Get("Content-Encoding"))
		}
	})
	t.Run("Brotli has higher priority than gzip", func(t *testing.T) {
		r := gin.New()
		r.Use(Compress())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "test")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept-Encoding", "gzip, br")
		r.ServeHTTP(w, req)

		if w.Header().Get("Content-Encoding") != "br" {
			t.Fatalf("Unsupported content encoding %s", w.Header().Get("Content-Encoding"))
		}
	})

}
