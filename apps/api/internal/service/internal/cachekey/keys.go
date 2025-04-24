// Package cachekey provides standardized cache key generation for the application
package cachekey

import (
	"fmt"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
)

// ForRepository generates a cache key for repository data
// Format: repo/{owner}--{repo}.{ext}
func ForRepository(r *model.Repository, ext string) string {
	return fmt.Sprintf("repo/%s--%s.%s", r.Owner, r.RepoName, ext)
}

// ForContributors generates a cache key for contributor data with versioning
// Format: contributors/v1.2/{owner}--{repo}.{ext}
func ForContributors(r *model.Repository, ext string) string {
	return fmt.Sprintf("contributors/v1.2/%s--%s.%s", r.Owner, r.RepoName, ext)
}

// ForImage generates a cache key for repository image
// Format: image/{owner}--{repo}--{anon|noanon}_{maxCount}_{columns}.{ext}
func ForImage(r *model.Repository, options *renderer.RendererOptions, ext string, includeAnonymous bool) string {
	anonStr := "noanon"
	if includeAnonymous {
		anonStr = "anon"
	}

	return fmt.Sprintf("image/%s--%s--%s_%d_%d.%s",
		r.Owner, r.RepoName, anonStr, options.MaxCount, options.Columns, ext)
}
