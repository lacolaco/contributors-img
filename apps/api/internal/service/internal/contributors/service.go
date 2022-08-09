package contributors

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/service/internal/cache"
	"contrib.rocks/apps/api/internal/service/logger"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/env"
	"contrib.rocks/libs/goutils/model"
	"github.com/google/go-github/v45/github"
)

type Service struct {
	env           env.Environment
	cacheService  *cache.Service
	githubClient  *github.Client
	loggingClient *logging.Client
}

func New(cfg *config.Config, c *cache.Service, gh *github.Client, l *logging.Client) *Service {
	return &Service{cfg.Env, c, gh, l}
}

func (s *Service) GetContributors(c context.Context, r *model.Repository) (*model.RepositoryContributors, error) {
	ctx, span := tracing.DefaultTracer.Start(c, "contributors.Service.GetContributors")
	defer span.End()
	log := logger.FromContext(ctx)

	cacheKey := createContributorsJSONCacheKey(r)
	// restore cache
	var cache *model.RepositoryContributors
	err := s.cacheService.GetJSON(ctx, cacheKey, &cache)
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
	err = s.cacheService.SaveJSON(ctx, cacheKey, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) sendCacheMissLog(c context.Context, key string) {
	s.loggingClient.Logger("contributors-json-cache-miss").Log(logging.Entry{
		Labels: map[string]string{
			"environment": string(s.env),
		},
		Payload: key,
	})
}

func createContributorsJSONCacheKey(r *model.Repository) string {
	// `contributors-json-cache/v1.2/${repository.owner}--${repository.repo}.json`;
	return fmt.Sprintf("contributors-json-cache/v1.2/%s--%s.json", r.Owner, r.RepoName)
}
