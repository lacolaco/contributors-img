package image

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api/go/dataurl"
	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
	"contrib.rocks/apps/api/go/util"
	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/tracing"
	"golang.org/x/sync/errgroup"
)

func New(cache appcache.AppCache) *Service {
	return &Service{cache}
}

type Service struct {
	cache appcache.AppCache
}

func (s *Service) GetImage(c context.Context, repo *model.Repository, options *renderer.RendererOptions, includeAnonymous bool) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "image.Service.GetImage")
	defer span.End()
	log := logger.LoggerFromContext(ctx)

	options = normalizeRendererOptions(options)
	cacheKey := createImageCacheKey(repo, options, "svg", includeAnonymous)

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

func (s *Service) RenderImage(c context.Context, data *model.RepositoryContributors, options *renderer.RendererOptions, includeAnonymous bool) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "image.Service.RenderImage")
	defer span.End()

	options = normalizeRendererOptions(options)
	cacheKey := createImageCacheKey(data.Repository, options, "svg", includeAnonymous)

	data, err := s.normalizeContributors(ctx, data, options, includeAnonymous)
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

func (s *Service) normalizeContributors(ctx context.Context, base *model.RepositoryContributors, options *renderer.RendererOptions, includeAnonymous bool) (*model.RepositoryContributors, error) {
	// filter by includeAnonymous
	contributors := make([]*model.Contributor, 0, len(base.Contributors))
	for _, c := range base.Contributors {
		if !includeAnonymous && c.ID == 0 {
			// skip anonymous contributors
			continue
		}
		contributors = append(contributors, c)
	}
	// get formatted data
	maxCount := util.Min(options.MaxCount, len(contributors))
	data := &model.RepositoryContributors{
		Repository:      base.Repository,
		StargazersCount: base.StargazersCount,
		Contributors:    make([]*model.Contributor, maxCount),
	}
	copy(data.Contributors, contributors)
	// convert avatar images to data urls
	eg, ctx := errgroup.WithContext(ctx)
	results := make([]string, len(data.Contributors))
	for i, c := range data.Contributors {
		index := i
		avatarURL := c.AvatarURL
		eg.Go(func() error {
			// ignore errors
			d, _ := dataurl.Convert(ctx, avatarURL, map[string]string{
				"size": fmt.Sprint(options.ItemSize),
				"s":    fmt.Sprint(options.ItemSize),
			})
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
