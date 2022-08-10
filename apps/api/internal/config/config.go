package config

import (
	"context"
	"fmt"
	"os"

	"contrib.rocks/libs/goutils/env"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port            string
	Env             env.Environment
	GitHubAuthToken string
	CacheBucketName string
	projectID       string
}

func (c *Config) ProjectID() string {
	return c.projectID
}

func Load() (*Config, error) {
	var config Config
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "3333"
	}
	config.Env = env.FromString(os.Getenv("APP_ENV"))
	config.GitHubAuthToken = os.Getenv("GITHUB_AUTH_TOKEN")
	if config.GitHubAuthToken == "" {
		return nil, fmt.Errorf("GITHUB_AUTH_TOKEN is required")
	}
	config.CacheBucketName = os.Getenv("CACHE_STORAGE_BUCKET")
	config.projectID = findProjectID()
	return &config, nil
}

func findProjectID() string {
	cred, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		return ""
	}
	return cred.ProjectID
}
