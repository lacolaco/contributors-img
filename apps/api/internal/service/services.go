package service

import (
	"io"

	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/service/internal/cache"
	"contrib.rocks/apps/api/internal/service/internal/contributors"
	"contrib.rocks/apps/api/internal/service/internal/image"
	"contrib.rocks/apps/api/internal/service/internal/usage"
	"contrib.rocks/apps/api/internal/service/logger"
	"contrib.rocks/libs/goutils/apiclient"
)

type ServicePack struct {
	ContributorsService ContributorsService
	UsageService        UsageService
	ImageService        ImageService
	LoggerFactory       logger.LoggerFactory
	DefaultLogger       logger.Logger

	closables []io.Closer
}

func NewServicePack(cfg *config.Config) *ServicePack {
	closables := []io.Closer{}
	storageClient := apiclient.NewStorageClient()
	loggingClient := apiclient.NewLoggingClient(cfg.ProjectID())
	closables = append(closables, loggingClient)
	githubClient := apiclient.NewGitHubClient(cfg.GitHubAuthToken)
	cacheService := cache.New(cfg, storageClient)
	loggingService := logger.NewLoggerFactory(cfg, loggingClient)

	return &ServicePack{
		ContributorsService: contributors.New(cfg, cacheService, githubClient, loggingClient),
		UsageService:        usage.New(cfg, loggingClient),
		ImageService:        image.New(cfg, cacheService, loggingClient),
		LoggerFactory:       loggingService,
		DefaultLogger:       loggingService.Logger("api/default"),

		closables: closables,
	}
}

func (sp *ServicePack) Close() {
	for _, fn := range sp.closables {
		fn.Close()
	}
}
