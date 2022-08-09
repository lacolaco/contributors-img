package config

import (
	"context"
	"fmt"

	"contrib.rocks/libs/goutils/env"
	"github.com/gobuffalo/envy"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port            string
	Env             env.Environment
	GitHubAuthToken string
	CacheBucketName string
	projectID       string // readonly
}

func (c *Config) ProjectID() string {
	return c.projectID
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
	config.projectID = getProjectID()
	return &config, nil
}

func getProjectID() string {
	cred, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(err)
	}
	if cred.ProjectID == "" {
		panic(fmt.Errorf("project id is not found"))
	}
	return cred.ProjectID
}
