package service

import (
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/service/contributors"
	"contrib.rocks/apps/api/internal/service/image"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/service/usage"
	"contrib.rocks/libs/goutils/apiclient"
)

type ServicePack struct {
	ContributorsService contributors.Service
	UsageService        usage.Service
	ImageService        *image.Service
}

func NewServicePack(cfg *config.Config) *ServicePack {
	gh := apiclient.NewGitHubClient(cfg.GitHubAuthToken)

	var cache appcache.AppCache
	if cfg.ProjectID() != "" && cfg.CacheBucketName != "" {
		storageClient := apiclient.NewStorageClient()
		cache = appcache.NewGCSCache(storageClient, cfg.CacheBucketName)
	} else {
		cache = appcache.NewMemoryCache()
	}

	return &ServicePack{
		ContributorsService: contributors.New(gh, cache),
		ImageService:        image.New(cache),
		UsageService:        usage.New(),
	}
}
