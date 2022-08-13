package app

import (
	"fmt"
	"net/http"

	"contrib.rocks/apps/api/internal/api"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/env"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %s", err.Error())
	}
	fmt.Printf("config: %+v\n", cfg)

	if cfg.Env == env.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	closeTracer := tracing.InitTraceProvider(cfg)
	defer closeTracer()
	loggerFactory := logger.NewLoggerFactory(cfg)
	defer loggerFactory.Close()
	sp := service.NewServicePack(cfg)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(config.Middleware(cfg))
	r.Use(tracing.Middleware(cfg))
	r.Use(logger.Middleware(loggerFactory))
	r.Use(logger.MiddlewareZap(cfg))
	r.Use(errorHandler())
	r.Use(requestLogger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	api.Setup(r, sp)

	fmt.Printf("Listening on http://localhost:%s\n", cfg.Port)
	return r.Run(fmt.Sprintf(":%s", cfg.Port))
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		logger.LoggerFromContext(c.Request.Context()).Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log := logger.LoggerFromContext(c.Request.Context())

		log.With(logger.LogGroupLabel("requestlog")).Sugar().Info(map[string]string{
			"status":    fmt.Sprintf("%d", c.Writer.Status()),
			"method":    c.Request.Method,
			"host":      c.Request.Host,
			"url":       c.Request.URL.String(),
			"referer":   c.Request.Referer(),
			"userAgent": c.Request.UserAgent(),
		})
	}
}
