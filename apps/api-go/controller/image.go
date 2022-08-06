package controller

import (
	"fmt"
	"net/http"

	"contrib.rocks/apps/api-go/service"
	"contrib.rocks/libs/goutils/renderer"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	cs *service.ContributorsService
	is *service.ImageService
	us *service.UsageService
}

func NewImageController(cs *service.ContributorsService, is *service.ImageService, us *service.UsageService) *ImageController {
	return &ImageController{cs, is, us}
}

func (c *ImageController) Get(ctx *gin.Context) {
	var params ImageParams
	if err := params.Bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("params=%+v\n", params)

	// get data
	data, err := c.cs.GetContributors(ctx, params.Repository.Object())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("data=%+v\n", data)
	// get image
	file, err := c.is.GetImage(ctx, &service.GetImageParams{
		Repository: params.Repository.Object(),
		RendererOptions: &renderer.RendererOptions{
			MaxCount: params.MaxCount,
			Columns:  params.Columns,
		},
		Data: data,
	})
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r := file.Reader()
	defer r.Close()
	ctx.DataFromReader(http.StatusOK, file.Size(), file.ContentType(), r, nil)
	// collect usage stats
	if err := c.us.Collect(ctx, data, params.Via); err != nil {
		fmt.Printf("error collecting usage: %s\n", err)
	}
}
