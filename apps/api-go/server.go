package main

import (
	"fmt"

	"contrib.rocks/apps/api-go/controller"
	"contrib.rocks/apps/api-go/core"
	"contrib.rocks/apps/api-go/middleware"
	"contrib.rocks/apps/api-go/service"
	"contrib.rocks/libs/goutils/config"
	"github.com/gin-gonic/gin"
)

func StartServer(port string) error {
	i := core.NewInfrastructure()
	defer i.Close()

	if i.Env == config.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(middleware.AppDefault()...)

	cs := service.NewContributorsService(i)
	is := service.NewImageService(i)
	usage := service.NewUsageService(i)
	{
		image := controller.NewImageController(cs, is, usage)
		r.GET("/image", image.Get)
	}

	fmt.Printf("Listening on http://localhost:%s\n", port)
	return r.Run(fmt.Sprintf(":%s", port))
}
