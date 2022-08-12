package logger

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	loggerContextKey contextKey = contextKey("logger")
)

func Middleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := wrapContext(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func wrapContext(c context.Context, logger Logger) context.Context {
	return context.WithValue(c, loggerContextKey, logger)
}

func FromContext(c context.Context) Logger {
	return c.Value(loggerContextKey).(Logger)
}
