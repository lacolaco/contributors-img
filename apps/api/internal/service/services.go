package service

import (
	"io"

	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/service/internal/cache"
	"contrib.rocks/apps/api/internal/service/internal/contributors"
	"contrib.rocks/apps/api/internal/service/internal/image"
	"contrib.rocks/apps/api/internal/service/internal/usage"
	"contrib.rocks/libs/goutils/apiclient"
)

type ServicePack struct {
	ContributorsService ContributorsService
	UsageService        UsageService
	ImageService        ImageService

	closables []io.Closer
}

func NewServicePack(cfg *config.Config) *ServicePack {
	closables := []io.Closer{}
	storageClient := apiclient.NewStorageClient()
	loggingClient := apiclient.NewLoggingClient()
	closables = append(closables, loggingClient)
	githubClient := apiclient.NewGitHubClient(cfg.GitHubAuthToken)
	cacheService := cache.New(cfg, storageClient)

	return &ServicePack{
		ContributorsService: contributors.New(cacheService, githubClient),
		UsageService:        usage.New(loggingClient, cfg),
		ImageService:        image.New(cacheService),
		closables:           closables,
	}
}

func (sp *ServicePack) Close() {
	for _, fn := range sp.closables {
		fn.Close()
	}
}
