package middleware

import (
	"contrib.rocks/libs/goutils/config"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func AppDefault() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		compression(),
		environment(),
	}
}

func environment() gin.HandlerFunc {
	env := config.GetEnv()
	return func(ctx *gin.Context) {
		ctx.Set(config.EnvContextKey, env)
	}
}

func compression() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}
