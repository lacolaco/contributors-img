package config

import (
	"os"
	"testing"

	"contrib.rocks/libs/go/env"
	"github.com/joho/godotenv"
)

func prepareEnv(t *testing.T) {
	os.Clearenv()
	err := godotenv.Load("../testing/.env")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Clearenv()
	})
}

func TestConfig_Load(t *testing.T) {
	t.Run("PORT is 3333 by default", func(t *testing.T) {
		prepareEnv(t)
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.Port != "3333" {
			t.Fatalf("Expected port to be 3333, got %s", config.Port)
		}
	})
	t.Run("PORT is set", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("PORT", "9000")
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.Port != "9000" {
			t.Fatalf("Expected port to be 9000, got %s", config.Port)
		}
	})
	t.Run("GITHUB_AUTH_TOKEN is required", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("GITHUB_AUTH_TOKEN", "")
		config, err := Load()
		if err == nil {
			t.Fatalf("Expected error, got nil: %+v", config)
		}
	})
	t.Run("APP_ENV is development by default", func(t *testing.T) {
		prepareEnv(t)
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.Env != env.EnvDevelopment {
			t.Fatalf("Expected env to be development, got %s", config.Env)
		}
	})
	t.Run("APP_ENV is set", func(tt *testing.T) {
		prepareEnv(t)
		os.Setenv("APP_ENV", "staging")
		config, err := Load()
		if err != nil {
			tt.Fatal(err)
		}
		if config.Env != env.EnvStaging {
			tt.Fatalf("Expected env to be staging, got %s", config.Env)
		}
	})
	t.Run("CACHE_STORAGE_BUCKET is empty by default", func(t *testing.T) {
		prepareEnv(t)
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.CacheBucketName != "" {
			t.Fatalf("Expected cache bucket name to be empty, got %s", config.CacheBucketName)
		}
	})
	t.Run("CACHE_STORAGE_BUCKET is set", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("CACHE_STORAGE_BUCKET", "test")
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.CacheBucketName != "test" {
			t.Fatalf("Expected cache bucket name to be test, got %s", config.CacheBucketName)
		}
	})
	t.Run("ProjectID is empty by default", func(t *testing.T) {
		prepareEnv(t)
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.ProjectID() != "" {
			t.Fatalf("Expected project ID to be empty, got %s", config.ProjectID())
		}
	})
}
