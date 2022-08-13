package config

import (
	"contrib.rocks/libs/goutils/env"
	"golang.org/x/oauth2/google"
)

func NewTestConfig() *Config {
	return &Config{
		Port:            "3333",
		Env:             env.EnvDevelopment,
		GitHubAuthToken: "test",
		CacheBucketName: "",

		googleCredentials: &google.Credentials{
			ProjectID: "test",
		},
	}
}
