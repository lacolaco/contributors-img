package config

import (
	"fmt"

	"contrib.rocks/libs/goutils/env"
	"github.com/gobuffalo/envy"
)

type Config struct {
	Port            string
	Env             env.Environment
	GitHubAuthToken string
	CacheBucketName string
}

func Load() (*Config, error) {
	envy.Load()
	var config Config
	config.Port = envy.Get("PORT", "3333")
	config.Env = env.FromString(envy.Get("APP_ENV", "development"))
	config.GitHubAuthToken = envy.Get("GITHUB_AUTH_TOKEN", "")
	if config.GitHubAuthToken == "" {
		return nil, fmt.Errorf("GITHUB_AUTH_TOKEN is required")
	}
	config.CacheBucketName = envy.Get("CACHE_STORAGE_BUCKET", "")
	return &config, nil
}
