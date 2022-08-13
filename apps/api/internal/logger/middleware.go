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
	loggerFactoryContextKey
)

func Middleware(factory LoggerFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, loggerFactoryContextKey, factory)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func MiddlewareZap(cfg *config.Config) gin.HandlerFunc {
	baseLogger := buildBaseLogger(cfg)

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := baseLogger.WithOptions(
			withTracing(ctx),
		)
		ctx = context.WithValue(ctx, loggerContextKey, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		logger.Debug("debug test")
		logger.Info("info test")
	}
}

func LoggerFromContext(c context.Context) *zap.Logger {
	return c.Value(loggerContextKey).(*zap.Logger)
}

func LoggerFactoryFromContext(c context.Context) LoggerFactory {
	return c.Value(loggerFactoryContextKey).(LoggerFactory)
}
