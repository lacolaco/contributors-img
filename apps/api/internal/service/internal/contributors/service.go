package contributors

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"github.com/google/go-github/v45/github"
)

type Service struct {
	githubClient    *github.Client
	cache           appcache.AppCache
	cacheMissLogger logger.Logger
}

func New(gh *github.Client, cache appcache.AppCache, cacheMissLogger logger.Logger) *Service {
	return &Service{gh, cache, cacheMissLogger}
}

func (s *Service) GetContributors(c context.Context, r *model.Repository) (*model.RepositoryContributors, error) {
	ctx, span := tracing.DefaultTracer.Start(c, "contributors.Service.GetContributors")
	defer span.End()
	log := logger.FromContext(ctx)

	cacheKey := createContributorsJSONCacheKey(r)
	// restore cache
	var cache *model.RepositoryContributors
	err := s.cache.GetJSON(ctx, cacheKey, &cache)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		log.Debug(ctx, logger.NewEntry(fmt.Sprintf("restored contributors-json from cache: %s", cacheKey)))
		return cache, nil
	}
	s.sendCacheMissLog(ctx, cacheKey)
	// get contributors from github
	data, err := resolveRepositoryData(s.githubClient, ctx, r)
	if err != nil {
		return nil, err
	}
	// save cache
	err = s.cache.SaveJSON(ctx, cacheKey, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) sendCacheMissLog(c context.Context, key string) {
	s.cacheMissLogger.Log(c, logging.Entry{
		Payload: key,
	})
}

func createContributorsJSONCacheKey(r *model.Repository) string {
	// `contributors-json-cache/v1.2/${repository.owner}--${repository.repo}.json`;
	return fmt.Sprintf("contributors-json-cache/v1.2/%s--%s.json", r.Owner, r.RepoName)
}
