package service

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api-go/core"
	"contrib.rocks/apps/api-go/infrastructure"
	"contrib.rocks/libs/goutils/model"
)

type ContributorsService struct {
	cache *infrastructure.CacheStorage
	ghc   *GitHubService
}

func NewContributorsService(i *core.Infrastructure) *ContributorsService {
	return &ContributorsService{i.Cache, NewGitHubService(i)}
}

func (s *ContributorsService) GetContributors(ctx context.Context, r *model.Repository) (*model.RepositoryContributors, error) {
	cacheKey := createContributorsJSONCacheKey(r)
	// restore cache
	var cache *model.RepositoryContributors
	err := s.cache.GetJSON(ctx, cacheKey, cache)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		return cache, nil
	}
	data, err := s.ghc.GetContributors(ctx, r)
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

func createContributorsJSONCacheKey(r *model.Repository) string {
	// `contributors-json-cache/v1.2/${repository.owner}--${repository.repo}.json`;
	return fmt.Sprintf("contributors-json-cache/v1.2/%s--%s.json", r.Owner, r.RepoName)
}
