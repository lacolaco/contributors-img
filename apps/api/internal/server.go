package app

import (
	"context"
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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

	closeLogger := logger.InitLoggerFactory(cfg)
	defer closeLogger()

	sp := service.NewServicePack(cfg)
	defer sp.Close()

	tp := tracing.InitTraceProvider(cfg)
	defer tp.Shutdown(context.Background())

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(env.Middleware(cfg.Env))
	r.Use(logger.Middleware(sp.DefaultLogger))
	r.Use(errorHandler())
	r.Use(otelgin.Middleware("api", otelgin.WithTracerProvider(tp)))
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	api.Setup(r, sp)

	fmt.Printf("Listening on http://localhost:%s\n", cfg.Port)
	return r.Run(fmt.Sprintf(":%s", cfg.Port))
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log := logger.FromContext(c.Request.Context())
		err := c.Errors.Last()
		if err == nil {
			return
		}
		log.Error(c, logger.NewEntry(err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
	}
}
