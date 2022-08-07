package app

import (
	"fmt"
	"net/http"

	"contrib.rocks/apps/api-go/internal/api"
	"contrib.rocks/apps/api-go/internal/config"
	"contrib.rocks/apps/api-go/internal/service"
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

	sp := service.NewServicePack(cfg)
	defer sp.Close()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
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
