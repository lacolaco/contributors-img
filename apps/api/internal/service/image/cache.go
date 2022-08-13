package image

import (
	"fmt"

	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
)

func createImageCacheKey(r *model.Repository, options *renderer.RendererOptions, ext string) string {
	return fmt.Sprintf("image-cache/%s--%s--%d_%d.%s",
		r.Owner, r.RepoName, options.MaxCount, options.Columns, ext)
}
