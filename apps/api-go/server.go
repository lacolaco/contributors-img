package main

import (
	"fmt"

	"contrib.rocks/apps/api-go/controller"
	"contrib.rocks/apps/api-go/core"
	"contrib.rocks/apps/api-go/middleware"
	"contrib.rocks/apps/api-go/service"
	"contrib.rocks/libs/goutils"
	"github.com/gin-gonic/gin"
)

func StartServer(port string) error {
	env := goutils.GetEnv()
	if env == goutils.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.AppDefault()...)

	i := core.NewInfrastructure()
	is := service.NewImageService(i)
	{
		image := controller.NewImageController(is)
		r.GET("/image", image.Get)
	}

	fmt.Printf("Listening on http://localhost:%s\n", port)
	return r.Run(fmt.Sprintf(":%s", port))
}
