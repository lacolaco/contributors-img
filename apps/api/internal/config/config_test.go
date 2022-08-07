package config

import (
	"testing"

	"contrib.rocks/libs/goutils/env"
	"github.com/gobuffalo/envy"
)

func TestConfig_Load(t *testing.T) {
	t.Run("PORT is 3333 by default", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "test")
			config, err := Load()
			if err != nil {
				tt.Fatal(err)
			}
			if config.Port != "3333" {
				tt.Fatalf("Expected port to be 3333, got %s", config.Port)
			}
		})
	})
	t.Run("PORT is set", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "test")
			envy.Set("PORT", "9000")
			config, err := Load()
			if err != nil {
				tt.Fatal(err)
			}
			if config.Port != "9000" {
				tt.Fatalf("Expected port to be 9000, got %s", config.Port)
			}
		})
	})
	t.Run("GITHUB_AUTH_TOKEN is required", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "")
			config, err := Load()
			if err == nil {
				tt.Fatalf("Expected error, got nil: %+v", config)
			}
		})
	})
	t.Run("APP_ENV is development by default", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "test")
			config, err := Load()
			if err != nil {
				tt.Fatal(err)
			}
			if config.Env != env.EnvDevelopment {
				tt.Fatalf("Expected env to be development, got %s", config.Env)
			}
		})
	})
	t.Run("APP_ENV is set", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "test")
			envy.Set("APP_ENV", "staging")
			config, err := Load()
			if err != nil {
				tt.Fatal(err)
			}
			if config.Env != env.EnvStaging {
				tt.Fatalf("Expected env to be staging, got %s", config.Env)
			}
		})
	})
	t.Run("CACHE_STORAGE_BUCKET is empty by default", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "test")
			config, err := Load()
			if err != nil {
				tt.Fatal(err)
			}
			if config.CacheBucketName != "" {
				tt.Fatalf("Expected cache bucket name to be empty, got %s", config.CacheBucketName)
			}
		})
	})
	t.Run("CACHE_STORAGE_BUCKET is set", func(tt *testing.T) {
		envy.Temp(func() {
			envy.Set("GITHUB_AUTH_TOKEN", "test")
			envy.Set("CACHE_STORAGE_BUCKET", "test")
			config, err := Load()
			if err != nil {
				tt.Fatal(err)
			}
			if config.CacheBucketName != "test" {
				tt.Fatalf("Expected cache bucket name to be test, got %s", config.CacheBucketName)
			}
		})
	})
}
