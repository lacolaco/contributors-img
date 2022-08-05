package renderer

import (
	"bytes"
	"io"

	"contrib.rocks/libs/goutils/model"
)

type RendererOptions struct {
	MaxCount int
	Columns  int
}

type renderer struct {
	options *RendererOptions
}

func NewRenderer(o *RendererOptions) *renderer {
	return &renderer{o}
}

var _ Image = &svgImage{}

type svgImage []byte

func (i svgImage) Reader() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(i))
}
func (i svgImage) Size() int64 {
	return int64(len(i))
}
func (i svgImage) Bytes() []byte {
	return i
}
func (i svgImage) ContentType() string {
	return "image/svg+xml"
}

func (r *renderer) Render(contributors *model.RepositoryContributors) (Image, error) {
	return svgImage([]byte("<svg></svg>")), nil
}
