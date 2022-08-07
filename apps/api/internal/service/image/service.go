package image

import (
	"context"
	"fmt"
	"sync"

	"contrib.rocks/apps/api/internal/service/cache"
	"contrib.rocks/libs/goutils"
	"contrib.rocks/libs/goutils/dataurl"
	"contrib.rocks/libs/goutils/model"
	"contrib.rocks/libs/goutils/renderer"
)

type Service struct {
	cacheService *cache.Service
}

func New(c *cache.Service) *Service {
	return &Service{c}
}

type GetImageParams struct {
	Repository      *model.Repository
	RendererOptions *renderer.RendererOptions
	Data            *model.RepositoryContributors
}

func (s *Service) GetImage(ctx context.Context, p *GetImageParams) (model.FileHandle, error) {
	// set default options
	const (
		defaultMaxCount = 100
		defaultColumns  = 12
		defaultItemSize = 64
	)
	if p.RendererOptions.MaxCount < 1 {
		p.RendererOptions.MaxCount = defaultMaxCount
	}
	if p.RendererOptions.Columns < 1 {
		p.RendererOptions.Columns = defaultColumns
	}
	p.RendererOptions.ItemSize = defaultItemSize

	cacheKey := createImageCacheKey(p, "svg")
	// restore cached image
	cache, err := s.restoreCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		fmt.Printf("GetImage: restored from cache: %s\n", cacheKey)
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

func (s *Service) restoreCache(ctx context.Context, key string) (model.FileHandle, error) {
	cache, err := s.cacheService.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func (s *Service) render(ctx context.Context, p *GetImageParams) (renderer.Image, error) {
	// get data
	maxCount := goutils.Min(p.RendererOptions.MaxCount, len(p.Data.Contributors))
	data := &model.RepositoryContributors{
		Repository:      p.Data.Repository,
		StargazersCount: p.Data.StargazersCount,
		Contributors:    make([]*model.Contributor, maxCount),
	}
	copy(data.Contributors, p.Data.Contributors)
	// convert avatar images to data urls
	var wg sync.WaitGroup
	chs := make([]chan string, 0, len(data.Contributors))
	for i, c := range data.Contributors {
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
		}(chs[i], c.AvatarURL, p.RendererOptions.ItemSize)
	}
	wg.Wait()
	for i, ch := range chs {
		data.Contributors[i].AvatarURL = <-ch
	}
	// render image
	r := renderer.NewRenderer(p.RendererOptions)
	image := r.Render(data)
	return image, nil
}

func (s *Service) saveCache(ctx context.Context, key string, image renderer.Image) error {
	return s.cacheService.Save(ctx, key, image.Bytes(), image.ContentType())
}

func createImageCacheKey(p *GetImageParams, ext string) string {
	// `image-cache/${repository.owner}--${repository.repo}--${params.maxCount}_${params.maxColumns}.${ext}`;
	return fmt.Sprintf("image-cache/%s--%s--%d_%d.%s",
		p.Repository.Owner, p.Repository.RepoName, p.RendererOptions.MaxCount, p.RendererOptions.Columns, ext)
}
