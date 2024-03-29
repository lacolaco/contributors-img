package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

type compressHandler struct {
	gzPool sync.Pool
	brPool sync.Pool
}

func Compress() gin.HandlerFunc {
	return (&compressHandler{
		gzPool: sync.Pool{
			New: func() any {
				w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
				return w
			},
		},
		brPool: sync.Pool{
			New: func() any {
				return brotli.NewWriterLevel(io.Discard, brotli.DefaultCompression)
			},
		},
	}).Handle
}

func (h *compressHandler) Handle(c *gin.Context) {
	if c.GetHeader("Accept-Encoding") == "" ||
		strings.Contains(c.GetHeader("Connection"), "Upgrade") ||
		strings.Contains(c.GetHeader("Accept"), "text/event-stream") {
		// skip compression
		return
	}

	if strings.Contains(c.GetHeader("Accept-Encoding"), "br") {
		h.handleBrotli(c)
	} else if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		h.handleGzip(c)
	}
}

func (h *compressHandler) handleBrotli(c *gin.Context) {
	ctx, span := otel.Tracer("compress").Start(c.Request.Context(), "compress.handleBrotli")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	w := h.brPool.Get().(*brotli.Writer)
	defer h.brPool.Put(w)
	defer w.Reset(io.Discard)
	w.Reset(c.Writer)

	c.Header("Content-Encoding", "br")
	c.Header("Vary", "Accept-Encoding")
	c.Writer = &brotliWriter{c.Writer, w}
	defer func() {
		w.Close()
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}

func (h *compressHandler) handleGzip(c *gin.Context) {
	ctx, span := otel.Tracer("compress").Start(c.Request.Context(), "compress.handleGzip")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	w := h.gzPool.Get().(*gzip.Writer)
	defer h.gzPool.Put(w)
	defer w.Reset(io.Discard)
	w.Reset(c.Writer)

	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")
	c.Writer = &gzipWriter{c.Writer, w}
	defer func() {
		w.Close()
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}

type brotliWriter struct {
	gin.ResponseWriter
	writer *brotli.Writer
}

func (w *brotliWriter) WriteString(s string) (int, error) {
	w.Header().Del("Content-Length")
	return w.writer.Write([]byte(s))
}

func (w *brotliWriter) Write(data []byte) (int, error) {
	w.Header().Del("Content-Length")
	return w.writer.Write(data)
}

func (w *brotliWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipWriter) WriteString(s string) (int, error) {
	w.Header().Del("Content-Length")
	return w.writer.Write([]byte(s))
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	w.Header().Del("Content-Length")
	return w.writer.Write(data)
}

func (w *gzipWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}
