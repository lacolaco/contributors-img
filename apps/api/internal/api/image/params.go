package image

import (
	"strings"

	"contrib.rocks/libs/go/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

type GetImageParams struct {
	Repository       model.RepositoryString `form:"repo" binding:"required"`
	MaxCount         int                    `form:"max"`
	Columns          int                    `form:"columns"`
	IncludeAnonymous bool                   `form:"anon"`
	Preview          bool                   `form:"preview"`
	Via              string
}

// MarshalLogObject implements zapcore.ObjectMarshaler
func (p GetImageParams) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("repository", string(p.Repository))
	enc.AddInt("max", p.MaxCount)
	enc.AddInt("columns", p.Columns)
	enc.AddBool("anon", p.IncludeAnonymous)
	enc.AddBool("preview", p.Preview)
	enc.AddString("via", p.Via)
	return nil
}

func (p *GetImageParams) bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindQuery(p); err != nil {
		return err
	}
	// validate repository name format
	if err := model.ValidateRepositoryName(string(p.Repository)); err != nil {
		return err
	}
	p.Via = "unknown"
	if strings.Contains(ctx.Request.UserAgent(), "github") {
		p.Via = "github"
	} else if strings.HasSuffix(ctx.Request.Host, "contrib.rocks") {
		p.Via = "preview"
	}
	return nil
}
