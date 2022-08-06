package config

import (
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
	var err error
	var config Config
	config.Port = envy.Get("PORT", "3333")
	config.Env = env.FromString(envy.Get("APP_ENV", "development"))
	config.GitHubAuthToken, err = envy.MustGet("GITHUB_AUTH_TOKEN")
	if err != nil {
		return nil, err
	}
	config.CacheBucketName = envy.Get("CACHE_STORAGE_BUCKET", "")
	return &config, nil
}
