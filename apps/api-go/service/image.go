package service

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api-go/core"
	"contrib.rocks/apps/api-go/infrastructure"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
)

type ImageService struct {
	cache *infrastructure.CacheStorage
	cs    *ContributorsService
}

func NewImageService(i *core.Infrastructure) *ImageService {
	return &ImageService{i.Cache, NewContributorsService(i)}
}

type GetImageParams struct {
	Repository      *model.Repository
	RendererOptions *renderer.RendererOptions
}

func (s *ImageService) GetImage(ctx context.Context, p *GetImageParams) (model.FileHandle, error) {
	cacheKey := createImageCacheKey(p, "svg")
	// restore cached image
	cache, err := s.restoreCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		return cache, nil
	}
	// render image
	image, err := s.render(ctx, p)
	if err != nil {
		return nil, err
	}
	// save image
	err = s.saveCache(ctx, cacheKey, image)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (s *ImageService) restoreCache(ctx context.Context, key string) (model.FileHandle, error) {
	cache, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func (s *ImageService) render(ctx context.Context, p *GetImageParams) (renderer.Image, error) {
	// get contributors
	contributors, err := s.cs.GetContributors(ctx, p.Repository)
	if err != nil {
		return nil, err
	}
	// render image
	r := renderer.NewRenderer(p.RendererOptions)
	return r.Render(contributors)
}

func (s *ImageService) saveCache(ctx context.Context, key string, image renderer.Image) error {
	return s.cache.Save(ctx, key, image.Bytes(), image.ContentType())
}

func createImageCacheKey(p *GetImageParams, ext string) string {
	// `image-cache/${repository.owner}--${repository.repo}--${params.maxCount}_${params.maxColumns}.${ext}`;
	return fmt.Sprintf("image-cache/%s--%s--%d_%d.%s",
		p.Repository.Owner, p.Repository.RepoName, p.RendererOptions.MaxCount, p.RendererOptions.Columns, ext)
}
