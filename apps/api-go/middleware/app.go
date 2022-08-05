package middleware

import (
	"contrib.rocks/libs/goutils"
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
	env := goutils.GetEnv()
	return func(ctx *gin.Context) {
		ctx.Set(goutils.EnvContextKey, env)
	}
}

func compression() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}
