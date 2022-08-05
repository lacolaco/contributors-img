package core

import "contrib.rocks/apps/api-go/infrastructure"

type Infrastructure struct {
	Cache        *infrastructure.CacheStorage
	GitHubClient *infrastructure.GitHubClient
}

func NewInfrastructure() *Infrastructure {
	return &Infrastructure{
		infrastructure.NewCacheStorage(),
		infrastructure.NewGitHubClient(),
	}
}
