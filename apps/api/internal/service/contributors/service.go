package contributors

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api/internal/service/cache"
	"contrib.rocks/libs/goutils/model"
	"github.com/google/go-github/v45/github"
)

type Service struct {
	cacheService *cache.Service
	githubClient *github.Client
}

func New(c *cache.Service, gh *github.Client) *Service {
	return &Service{c, gh}
}

func (s *Service) GetContributors(ctx context.Context, r *model.Repository) (*model.RepositoryContributors, error) {
	cacheKey := createContributorsJSONCacheKey(r)
	// restore cache
	var cache *model.RepositoryContributors
	err := s.cacheService.GetJSON(ctx, cacheKey, cache)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		fmt.Printf("GetContributors: restored from cache: %s\n", cacheKey)
		return cache, nil
	}
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

func createContributorsJSONCacheKey(r *model.Repository) string {
	// `contributors-json-cache/v1.2/${repository.owner}--${repository.repo}.json`;
	return fmt.Sprintf("contributors-json-cache/v1.2/%s--%s.json", r.Owner, r.RepoName)
}
