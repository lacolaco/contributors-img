package config

import "contrib.rocks/libs/goutils/env"

func NewTestConfig() *Config {
	return &Config{
		Port:            "3333",
		Env:             env.EnvDevelopment,
		GitHubAuthToken: "test",
		CacheBucketName: "",
		projectID:       "test",
	}
}
