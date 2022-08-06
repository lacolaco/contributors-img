package service

import (
	"io"

	"contrib.rocks/apps/api-go/internal/config"
	"contrib.rocks/apps/api-go/internal/service/cache"
	"contrib.rocks/apps/api-go/internal/service/contributors"
	"contrib.rocks/apps/api-go/internal/service/image"
	"contrib.rocks/apps/api-go/internal/service/usage"
	"contrib.rocks/libs/goutils/apiclient"
)

// TODO: make services less coupled (dependency injection?)
type ServicePack struct {
	CacheService        *cache.Service
	ContributorsService *contributors.Service
	UsageService        *usage.Service
	ImageService        *image.Service
	closables           []io.Closer
}

func NewServicePack(cfg *config.Config) *ServicePack {
	closables := []io.Closer{}
	storageClient := apiclient.NewStorageClient()
	loggingClient := apiclient.NewLoggingClient()
	closables = append(closables, loggingClient)
	githubClient := apiclient.NewGitHubClient(cfg.GitHubAuthToken)
	cacheService := cache.New(cfg, storageClient)

	return &ServicePack{
		CacheService:        cacheService,
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
