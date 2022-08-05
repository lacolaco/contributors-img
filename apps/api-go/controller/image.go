package controller

import (
	"fmt"
	"net/http"

	"contrib.rocks/apps/api-go/service"
	"contrib.rocks/libs/goutils/renderer"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	is *service.ImageService
}

func NewImageController(is *service.ImageService) *ImageController {
	return &ImageController{is}
}

func (c *ImageController) Get(ctx *gin.Context) {
	var params ImageParams
	if err := params.Bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%#v\n", params)

	file, err := c.is.GetImage(ctx, &service.GetImageParams{
		Repository: params.Repository.Object(),
		RendererOptions: &renderer.RendererOptions{
			MaxCount: params.MaxCount,
			Columns:  params.Columns,
		}})
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r := file.Reader()
	defer r.Close()
	ctx.DataFromReader(http.StatusOK, file.Size(), file.ContentType(), r, nil)

	// collect usage stats
}
