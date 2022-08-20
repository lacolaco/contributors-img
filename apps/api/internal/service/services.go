package service

import (
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/service/contributors"
	"contrib.rocks/apps/api/internal/service/image"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/service/usage"
	"contrib.rocks/libs/goutils/apiclient"
	"contrib.rocks/libs/goutils/github"
)

type ServicePack struct {
	ContributorsService *contributors.Service
	UsageService        *usage.Service
	ImageService        *image.Service
}

func NewServicePack(cfg *config.Config) *ServicePack {
	ghProvider := github.NewProvider(cfg.GitHubAuthToken)

	var cache appcache.AppCache
	if cfg.GoogleCredentials() != nil && cfg.CacheBucketName != "" {
		storageClient := apiclient.NewStorageClient()
		cache = appcache.NewGCSCache(storageClient, cfg.CacheBucketName)
	} else {
		cache = appcache.NewMemoryCache()
	}

	return &ServicePack{
		ContributorsService: contributors.New(ghProvider, cache),
		ImageService:        image.New(cache),
		UsageService:        usage.New(),
	}
}
