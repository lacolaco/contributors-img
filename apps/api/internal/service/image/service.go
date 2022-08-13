package image

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils"
	"contrib.rocks/libs/goutils/dataurl"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
	"golang.org/x/sync/errgroup"
)

func New(cache appcache.AppCache) *Service {
	return &Service{cache}
}

type Service struct {
	cache appcache.AppCache
}

func (s *Service) GetImage(c context.Context, repo *model.Repository, options *renderer.RendererOptions) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "image.Service.GetImage")
	defer span.End()
	log := logger.LoggerFromContext(ctx)

	options = normalizeRendererOptions(options)
	cacheKey := createImageCacheKey(repo, options, "svg")

	cache, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if cache == nil {
		s.sendCacheMissLog(ctx, cacheKey)
		return nil, nil
	}
	log.Debug(fmt.Sprintf("restored image from cache: %s", cacheKey))
	return cache, nil
}

func (s *Service) RenderImage(c context.Context, data *model.RepositoryContributors, options *renderer.RendererOptions) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "image.Service.RenderImage")
	defer span.End()

	options = normalizeRendererOptions(options)
	cacheKey := createImageCacheKey(data.Repository, options, "svg")

	data, err := s.normalizeContributors(ctx, data, options)
	if err != nil {
		return nil, err
	}

	image := renderer.NewRenderer(options).Render(data)

	err = s.cache.Save(c, cacheKey, image.Bytes(), image.ContentType())
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (s *Service) normalizeContributors(ctx context.Context, base *model.RepositoryContributors, options *renderer.RendererOptions) (*model.RepositoryContributors, error) {
	// get formatted data
	maxCount := goutils.Min(options.MaxCount, len(base.Contributors))
	data := &model.RepositoryContributors{
		Repository:      base.Repository,
		StargazersCount: base.StargazersCount,
		Contributors:    make([]*model.Contributor, maxCount),
	}
	copy(data.Contributors, base.Contributors)
	// convert avatar images to data urls
	eg, ctx := errgroup.WithContext(ctx)
	results := make([]string, len(data.Contributors))
	for i, c := range data.Contributors {
		index := i
		avatarURL := c.AvatarURL
		eg.Go(func() error {
			// ignore errors
			d, _ := dataurl.ResolveImageDataURL(ctx, avatarURL, options.ItemSize)
			results[index] = d
			return nil
		})
	}
	eg.Wait()
	for i, result := range results {
		data.Contributors[i].AvatarURL = result
	}
	return data, nil
}

func (s *Service) sendCacheMissLog(c context.Context, key string) {
	logger.LoggerFromContext(c).With(logger.LogGroup("image-cache-miss")).Info(
		fmt.Sprintf("image-cache-miss: %s", key),
	)
}
