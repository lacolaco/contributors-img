package core

import (
	"contrib.rocks/apps/api-go/infrastructure"
	"contrib.rocks/libs/goutils/config"
)

type Infrastructure struct {
	Env           config.Environment
	Cache         *infrastructure.CacheStorage
	GitHubClient  *infrastructure.GitHubClient
	LoggingClient *infrastructure.LoggingClient
}

func NewInfrastructure() *Infrastructure {
	return &Infrastructure{
		config.GetEnv(),
		infrastructure.NewCacheStorage(),
		infrastructure.NewGitHubClient(),
		infrastructure.NewLoggingClient(),
	}
}

func (i *Infrastructure) Close() {
	i.LoggingClient.Close()
}
