package service

import (
	"context"

	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
)

type UsageService interface {
	CollectUsage(ctx context.Context, r *model.RepositoryContributors, via string) error
}

type ImageService interface {
	GetImage(ctx context.Context, r *model.RepositoryContributors, options *renderer.RendererOptions) (model.FileHandle, error)
}
