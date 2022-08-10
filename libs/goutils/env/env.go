package env

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Environment string

const (
	EnvDevelopment = Environment("development")
	EnvStaging     = Environment("staging")
	EnvProduction  = Environment("production")
)

type contextKey string

const (
	envContextKey = contextKey("env")
)

// FromString app environment from string
func FromString(str string) Environment {
	switch str {
	case "production":
		return EnvProduction
	case "staging":
		return EnvStaging
	default:
		return EnvDevelopment
	}
}

func FromContext(ctx context.Context) Environment {
	return FromString(ctx.Value(envContextKey).(string))
}

func (e Environment) ContextWithEnvironment(parent context.Context) context.Context {
	return context.WithValue(parent, envContextKey, string(e))
}

func Middleware(e Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := e.ContextWithEnvironment(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
