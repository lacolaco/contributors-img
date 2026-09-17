package config

import (
	"os"
	"testing"

	"contrib.rocks/apps/api/go/env"
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
	t.Run("error when neither GitHub App nor GITHUB_AUTH_TOKEN is configured", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("GITHUB_AUTH_TOKEN", "")
		config, err := Load()
		if err == nil {
			t.Fatalf("Expected error, got nil: %+v", config)
		}
	})
	t.Run("falls back to GITHUB_AUTH_TOKEN when GitHub App is not configured", func(t *testing.T) {
		prepareEnv(t)
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.GitHubAuthToken != "test" {
			t.Fatalf("Expected GitHubAuthToken to be test, got %s", config.GitHubAuthToken)
		}
		if config.GitHubApp != nil {
			t.Fatalf("Expected GitHubApp to be nil, got %+v", config.GitHubApp)
		}
	})
	t.Run("GitHub App is configured when all three variables are set", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("GITHUB_APP_ID", "123")
		os.Setenv("GITHUB_APP_INSTALLATION_ID", "456")
		os.Setenv("GITHUB_APP_PRIVATE_KEY", "test-key")
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.GitHubApp == nil {
			t.Fatal("Expected GitHubApp to be set, got nil")
		}
		if config.GitHubApp.AppID != 123 {
			t.Fatalf("Expected AppID to be 123, got %d", config.GitHubApp.AppID)
		}
		if config.GitHubApp.InstallationID != 456 {
			t.Fatalf("Expected InstallationID to be 456, got %d", config.GitHubApp.InstallationID)
		}
		if config.GitHubApp.PrivateKey != "test-key" {
			t.Fatalf("Expected PrivateKey to be test-key, got %s", config.GitHubApp.PrivateKey)
		}
	})
	t.Run("GitHub App takes priority over GITHUB_AUTH_TOKEN when both are set", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("GITHUB_APP_ID", "123")
		os.Setenv("GITHUB_APP_INSTALLATION_ID", "456")
		os.Setenv("GITHUB_APP_PRIVATE_KEY", "test-key")
		config, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if config.GitHubApp == nil {
			t.Fatal("Expected GitHubApp to be set, got nil")
		}
	})
	t.Run("error when only some GitHub App variables are set", func(t *testing.T) {
		testCases := []struct {
			name           string
			appID          string
			installationID string
			privateKey     string
		}{
			{"only GITHUB_APP_ID", "123", "", ""},
			{"only GITHUB_APP_INSTALLATION_ID", "", "456", ""},
			{"only GITHUB_APP_PRIVATE_KEY", "", "", "test-key"},
			{"missing GITHUB_APP_PRIVATE_KEY", "123", "456", ""},
			{"missing GITHUB_APP_INSTALLATION_ID", "123", "", "test-key"},
			{"missing GITHUB_APP_ID", "", "456", "test-key"},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				prepareEnv(t)
				os.Setenv("GITHUB_APP_ID", tc.appID)
				os.Setenv("GITHUB_APP_INSTALLATION_ID", tc.installationID)
				os.Setenv("GITHUB_APP_PRIVATE_KEY", tc.privateKey)
				config, err := Load()
				if err == nil {
					t.Fatalf("Expected error, got nil: %+v", config)
				}
			})
		}
	})
	t.Run("error when GITHUB_APP_ID is not numeric", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("GITHUB_APP_ID", "not-a-number")
		os.Setenv("GITHUB_APP_INSTALLATION_ID", "456")
		os.Setenv("GITHUB_APP_PRIVATE_KEY", "test-key")
		config, err := Load()
		if err == nil {
			t.Fatalf("Expected error, got nil: %+v", config)
		}
	})
	t.Run("error when GITHUB_APP_INSTALLATION_ID is not numeric", func(t *testing.T) {
		prepareEnv(t)
		os.Setenv("GITHUB_APP_ID", "123")
		os.Setenv("GITHUB_APP_INSTALLATION_ID", "not-a-number")
		os.Setenv("GITHUB_APP_PRIVATE_KEY", "test-key")
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
