package image

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils"
	"contrib.rocks/libs/goutils/dataurl"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
)

type Service interface {
	GetImage(ctx context.Context, r *model.RepositoryContributors, options *renderer.RendererOptions) (model.FileHandle, error)
}

func New(cache appcache.AppCache) Service {
	return &serviceImpl{cache}
}

var _ Service = &serviceImpl{}

type serviceImpl struct {
	cache appcache.AppCache
}

type GetImageParams struct {
	RendererOptions *renderer.RendererOptions
	Data            *model.RepositoryContributors
}

func (s *serviceImpl) GetImage(c context.Context, data *model.RepositoryContributors, options *renderer.RendererOptions) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "image.Service.GetImage")
	defer span.End()
	log := logger.LoggerFromContext(ctx)

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

	cacheKey := createImageCacheKey(data.Repository, options, "svg")
	// restore cached image
	cache, err := s.restoreCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		log.Debug(ctx, logger.NewEntry(fmt.Sprintf("restored image from cache: %s", cacheKey)))
		return cache, nil
	}
	s.sendCacheMissLog(ctx, cacheKey)
	// render image
	image, err := s.render(ctx, data, options)
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

func (s *serviceImpl) restoreCache(ctx context.Context, key string) (model.FileHandle, error) {
	cache, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func (s *serviceImpl) render(c context.Context, data *model.RepositoryContributors, options *renderer.RendererOptions) (renderer.Image, error) {
	ctx, span := tracing.Tracer().Start(c, "image.Service.render")
	defer span.End()

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
			d, err := dataurl.ResolveImageDataURL(ctx, avatarUrl, itemSize)
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

func (s *serviceImpl) saveCache(c context.Context, key string, image renderer.Image) error {
	return s.cache.Save(c, key, image.Bytes(), image.ContentType())
}

func (s *serviceImpl) sendCacheMissLog(c context.Context, key string) {
	logger.LoggerFactoryFromContext(c).Logger("image-cache-miss").Log(c, logging.Entry{
		Payload: key,
	})
}

func createImageCacheKey(r *model.Repository, options *renderer.RendererOptions, ext string) string {
	// `image-cache/${repository.owner}--${repository.repo}--${params.maxCount}_${params.maxColumns}.${ext}`;
	return fmt.Sprintf("image-cache/%s--%s--%d_%d.%s",
		r.Owner, r.RepoName, options.MaxCount, options.Columns, ext)
}
