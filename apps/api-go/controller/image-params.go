package controller

import (
	"strings"

	"contrib.rocks/libs/goutils/model"
	"github.com/gin-gonic/gin"
)

type ImageParams struct {
	Repository model.RepositoryString `form:"repo" binding:"required"`
	MaxCount   int                    `form:"max"`
	Columns    int                    `form:"columns"`
	Preview    bool                   `form:"preview"`
	Via        string
}

const (
	defaultMaxCount = 100
	defaultColumns  = 12
)

func (p *ImageParams) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindQuery(p); err != nil {
		return err
	}
	// validate repository name format
	if err := model.ValidateRepositoryName(string(p.Repository)); err != nil {
		return err
	}
	// validate max count is upper than 1
	if p.MaxCount < 1 {
		p.MaxCount = defaultMaxCount
	}
	// validate columns is upper than 1
	if p.Columns < 1 {
		p.Columns = defaultColumns
	}
	p.Via = "unknown"
	if strings.Contains(ctx.Request.UserAgent(), "github") {
		p.Via = "github"
	} else if strings.HasSuffix(ctx.Request.Host, "contrib.rocks") {
		p.Via = "preview"
	}
	return nil
}
