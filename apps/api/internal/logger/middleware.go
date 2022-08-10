package logger

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
)

func Middleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.ContextWithLogger(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
		logger.Info(c.Request.Context(), logging.Entry{
			HTTPRequest: &logging.HTTPRequest{
				Request: c.Request,
			},
			Timestamp: time.Now(),
			Payload: map[string]string{
				"status":    fmt.Sprintf("%d", c.Writer.Status()),
				"method":    c.Request.Method,
				"host":      c.Request.Host,
				"url":       c.Request.URL.String(),
				"referer":   c.Request.Referer(),
				"userAgent": c.Request.UserAgent(),
			},
		})
	}
}

func FromContext(c context.Context) Logger {
	return fromContext(c)
}
