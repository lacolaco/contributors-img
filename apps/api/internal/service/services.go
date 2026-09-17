package service

import (
	"fmt"

	"contrib.rocks/apps/api/go/apiclient"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/github"
	"contrib.rocks/apps/api/internal/service/contributors"
	"contrib.rocks/apps/api/internal/service/image"
	"contrib.rocks/apps/api/internal/service/internal/appcache"
	"contrib.rocks/apps/api/internal/service/usage"
)

type ServicePack struct {
	ContributorsService *contributors.Service
	UsageService        *usage.Service
	ImageService        *image.Service
}

func NewServicePack(cfg *config.Config) (*ServicePack, error) {
	ghProvider, err := newGitHubProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub provider: %w", err)
	}

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
	}, nil
}

// newGitHubProvider prefers GitHub App installation authentication when
// configured, falling back to a static token for local development against
// an unregistered app. config.Load already rejects a partial App
// configuration and the case where neither is set.
func newGitHubProvider(cfg *config.Config) (contributors.GitHubClientProvider, error) {
	if cfg.GitHubApp != nil {
		return github.NewProvider(cfg.GitHubApp.AppID, cfg.GitHubApp.InstallationID, cfg.GitHubApp.PrivateKey)
	}
	return github.NewTokenProvider(cfg.GitHubAuthToken), nil
}
