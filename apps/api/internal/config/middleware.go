package config

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const configContextKey = contextKey("config")

func Middleware(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := wrapContext(c.Request.Context(), cfg)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func wrapContext(parent context.Context, cfg *Config) context.Context {
	return context.WithValue(parent, configContextKey, cfg)
}

func FromContext(ctx context.Context) *Config {
	return ctx.Value(configContextKey).(*Config)
}
