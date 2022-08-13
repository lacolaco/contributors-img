package logger

import (
	"context"

	"contrib.rocks/apps/api/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type contextKey int

const (
	loggerContextKey contextKey = iota
)

// Middleware returns a gin middleware that sets the logger in the context.
func Middleware(cfg *config.Config) gin.HandlerFunc {
	baseLogger := buildBaseLogger(cfg)

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := baseLogger.WithOptions(
			withTracing(ctx),
		)
		ctx = ContextWithLogger(ctx, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ContextWithLogger returns a new context with the given logger.
func ContextWithLogger(c context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(c, loggerContextKey, logger)
}

// LoggerFromContext returns the logger for the given context.
func LoggerFromContext(c context.Context) *zap.Logger {
	return c.Value(loggerContextKey).(*zap.Logger)
}
