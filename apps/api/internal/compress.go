package app

import (
	"compress/gzip"
	"fmt"
	"strings"

	"contrib.rocks/apps/api/internal/logger"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func compress() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.LoggerFromContext(c.Request.Context())

		if c.GetHeader("Accept-Encoding") == "" ||
			strings.Contains(c.GetHeader("Connection"), "Upgrade") ||
			strings.Contains(c.GetHeader("Accept"), "text/event-stream") {
			// skip compression
			return
		}

		if strings.Contains(c.GetHeader("Accept-Encoding"), "br") {
			log.Debug("[compress] compressing with brotli")
			w := brotli.NewWriterLevel(c.Writer, brotli.BestCompression)
			c.Header("Content-Encoding", "br")
			c.Header("Vary", "Accept-Encoding")
			c.Writer = &brotliWriter{c.Writer, w}
			defer func() {
				w.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
		} else if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			log.Debug("[compress] compressing with gzip")
			w, _ := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")
			c.Writer = &gzipWriter{c.Writer, w}
			defer func() {
				w.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
		} else {
			log.Info("[compress] unsupported encoding", zap.String("encoding", c.GetHeader("Accept-Encoding")))
		}
		c.Next()
	}
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
