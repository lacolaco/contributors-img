package image

import (
	"fmt"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
)

func createImageCacheKey(r *model.Repository, options *renderer.RendererOptions, ext string, includeAnonymous bool) string {
	var anonymous string
	if includeAnonymous {
		anonymous = "anon"
	} else {
		anonymous = "noanon"
	}
	return fmt.Sprintf("image-cache/%s--%s--%s_%d_%d.%s",
		r.Owner, r.RepoName, anonymous, options.MaxCount, options.Columns, ext)
}
