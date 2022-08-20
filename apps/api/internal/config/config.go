package config

import (
	"context"
	"fmt"
	"os"

	"contrib.rocks/libs/go/env"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port            string
	Env             env.Environment
	GitHubAuthToken string
	CacheBucketName string

	googleCredentials *google.Credentials
}

func (c *Config) GoogleCredentials() *google.Credentials {
	return c.googleCredentials
}

func (c *Config) ProjectID() string {
	if c.googleCredentials != nil {
		return c.googleCredentials.ProjectID
	}
	return ""
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
	config.googleCredentials = findGoogleCredentials()
	return &config, nil
}

func findGoogleCredentials() *google.Credentials {
	cred, _ := google.FindDefaultCredentials(context.Background())
	return cred
}
