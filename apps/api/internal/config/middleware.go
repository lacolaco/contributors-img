package config

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const configContextKey = contextKey("config")

func Middleware(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), configContextKey, cfg)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func FromContext(ctx context.Context) *Config {
	return ctx.Value(configContextKey).(*Config)
}
