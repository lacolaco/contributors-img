package logger

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	loggerContextKey        contextKey = contextKey("logger")
	loggerFactoryContextKey contextKey = contextKey("loggerFactory")
)

func Middleware(factory LoggerFactory, loggerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := factory.Logger(loggerName)
		ctx = context.WithValue(ctx, loggerFactoryContextKey, factory)
		ctx = context.WithValue(ctx, loggerContextKey, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func LoggerFromContext(c context.Context) Logger {
	return c.Value(loggerContextKey).(Logger)
}

func LoggerFactoryFromContext(c context.Context) LoggerFactory {
	return c.Value(loggerFactoryContextKey).(LoggerFactory)
}
