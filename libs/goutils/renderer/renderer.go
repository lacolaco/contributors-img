package renderer

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"math"

	"contrib.rocks/libs/goutils/model"
	svg "github.com/ajstarks/svgo"
)

type RendererOptions struct {
	MaxCount int
	Columns  int
	ItemSize int
}

type renderer struct {
	Options *RendererOptions
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
func (i svgImage) ContentType() string {
	return "image/svg+xml"
}
func (i svgImage) ETag() string {
	return fmt.Sprintf("%x", md5.Sum(i.Bytes()))
}
func (i svgImage) Bytes() []byte {
	return i
}

func (r *renderer) Render(data *model.RepositoryContributors) Image {
	w := bytes.NewBuffer(nil)
	r.buildSVG(w, data)
	return svgImage(w.Bytes())
}

func (r *renderer) buildSVG(w io.Writer, data *model.RepositoryContributors) {
	itemCount := len(data.Contributors)
	columns := math.Min(float64(r.Options.Columns), float64(itemCount))
	rows := math.Ceil(float64(itemCount) / float64(columns))
	gap := 4
	width := float64(r.Options.ItemSize)*columns + float64(gap)*(columns-1)
	height := float64(r.Options.ItemSize)*rows + float64(gap)*(rows-1)

	canvas := svg.New(w)
	canvas.Start(int(width), int(height), fmt.Sprintf(`viewBox="0 0 %d %d"`, int(width), int(height)))
	for i, c := range data.Contributors {
		x := (i % int(columns)) * (r.Options.ItemSize + gap)
		y := (i / int(columns)) * (r.Options.ItemSize + gap)
		// nested <svg> is not supported in svgo
		fmt.Fprintf(canvas.Writer, `<svg x="%d" y="%d" width="%d" height="%d">\n`, x, y, r.Options.ItemSize, r.Options.ItemSize)
		{
			fillId := fmt.Sprintf("fill%d", c.ID)
			canvas.Title(c.Login)
			canvas.Circle(r.Options.ItemSize/2, r.Options.ItemSize/2, r.Options.ItemSize/2, `stroke="#c0c0c0"`, `stroke-width="1"`, fmt.Sprintf(`fill="url(#%s)"`, fillId))
			canvas.Def()
			{
				canvas.Pattern(fillId, 0, 0, r.Options.ItemSize, r.Options.ItemSize, "user")
				{
					canvas.Image(0, 0, r.Options.ItemSize, r.Options.ItemSize, c.AvatarURL)
				}
				canvas.PatternEnd()
			}
			canvas.DefEnd()
		}
		canvas.End()
	}
	canvas.End()
}
