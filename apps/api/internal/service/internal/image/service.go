package image

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/service/internal/cache"
	"contrib.rocks/libs/goutils"
	"contrib.rocks/libs/goutils/dataurl"
	"contrib.rocks/libs/goutils/env"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
)

type Service struct {
	env           env.Environment
	cacheService  *cache.Service
	loggingClient *logging.Client
}

func New(cfg *config.Config, c *cache.Service, l *logging.Client) *Service {
	return &Service{cfg.Env, c, l}
}

type GetImageParams struct {
	RendererOptions *renderer.RendererOptions
	Data            *model.RepositoryContributors
}

func (s *Service) GetImage(ctx context.Context, r *model.RepositoryContributors, options *renderer.RendererOptions) (model.FileHandle, error) {
	// set default options
	const (
		defaultMaxCount = 100
		defaultColumns  = 12
		defaultItemSize = 64
	)
	if options.MaxCount < 1 {
		options.MaxCount = defaultMaxCount
	}
	if options.Columns < 1 {
		options.Columns = defaultColumns
	}
	options.ItemSize = defaultItemSize

	cacheKey := createImageCacheKey(r.Repository, options, "svg")
	// restore cached image
	cache, err := s.restoreCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		fmt.Printf("GetImage: restored from cache: %s\n", cacheKey)
		return cache, nil
	}
	s.sendCacheMissLog(ctx, cacheKey)
	// render image
	image, err := s.render(ctx, r, options)
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

func (s *Service) restoreCache(ctx context.Context, key string) (model.FileHandle, error) {
	cache, err := s.cacheService.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func (s *Service) render(ctx context.Context, data *model.RepositoryContributors, options *renderer.RendererOptions) (renderer.Image, error) {
	// get formatted data
	maxCount := goutils.Min(options.MaxCount, len(data.Contributors))
	formatted := &model.RepositoryContributors{
		Repository:      data.Repository,
		StargazersCount: data.StargazersCount,
		Contributors:    make([]*model.Contributor, maxCount),
	}
	copy(formatted.Contributors, data.Contributors)
	// convert avatar images to data urls
	var wg sync.WaitGroup
	chs := make([]chan string, 0, len(formatted.Contributors))
	for i, c := range formatted.Contributors {
		wg.Add(1)
		chs = append(chs, make(chan string, 1))
		go func(ret chan<- string, avatarUrl string, itemSize int) {
			defer wg.Done()
			d, err := dataurl.ResolveImageDataURL(avatarUrl, itemSize)
			if err != nil {
				ret <- ""
				return
			}
			ret <- d
		}(chs[i], c.AvatarURL, options.ItemSize)
	}
	wg.Wait()
	for i, ch := range chs {
		formatted.Contributors[i].AvatarURL = <-ch
	}
	// render image
	r := renderer.NewRenderer(options)
	image := r.Render(formatted)
	return image, nil
}

func (s *Service) saveCache(ctx context.Context, key string, image renderer.Image) error {
	return s.cacheService.Save(ctx, key, image.Bytes(), image.ContentType())
}

func (s *Service) sendCacheMissLog(ctx context.Context, key string) {
	s.loggingClient.Logger("image-cache-miss").Log(logging.Entry{
		Labels: map[string]string{
			"environment": string(s.env),
		},
		Payload: key,
	})
}

func createImageCacheKey(r *model.Repository, options *renderer.RendererOptions, ext string) string {
	// `image-cache/${repository.owner}--${repository.repo}--${params.maxCount}_${params.maxColumns}.${ext}`;
	return fmt.Sprintf("image-cache/%s--%s--%d_%d.%s",
		r.Owner, r.RepoName, options.MaxCount, options.Columns, ext)
}
