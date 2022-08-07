package app

import (
	"context"
	"fmt"
	"net/http"

	"contrib.rocks/apps/api/internal/api"
	"contrib.rocks/apps/api/internal/config"
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

	sp := service.NewServicePack(cfg)
	defer sp.Close()

	tp := tracing.InitTraceProvider(cfg)
	defer tp.Shutdown(context.Background())

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(otelgin.Middleware("api", otelgin.WithTracerProvider(tp)))
	r.Use(errorHandler)

	api.Setup(r, sp)

	fmt.Printf("Listening on http://localhost:%s\n", cfg.Port)
	return r.Run(fmt.Sprintf(":%s", cfg.Port))
}

func errorHandler(c *gin.Context) {
	c.Next()

	err := c.Errors.ByType(gin.ErrorTypePublic).Last()
	if err != nil {
		fmt.Fprint(gin.DefaultErrorWriter, err.Err)
		c.AbortWithError(http.StatusInternalServerError, err.Err)
	}
}
