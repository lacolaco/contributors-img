package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"contrib.rocks/apps/api/go/env"
	"golang.org/x/oauth2/google"
)

// GitHubAppConfig holds the credentials needed to authenticate as a GitHub App
// installation. All three fields must be set together; see resolveGitHubApp.
type GitHubAppConfig struct {
	AppID          int64
	InstallationID int64
	PrivateKey     string
}

type Config struct {
	Port            string
	Env             env.Environment
	GitHubAuthToken string
	GitHubApp       *GitHubAppConfig
	CacheBucketName string

	googleCredentials *google.Credentials
}

// String leaves the token out. StartServer prints the Config at startup, and
// that output is kept in Cloud Logging.
func (c *Config) String() string {
	return fmt.Sprintf("{Port:%s Env:%s CacheBucketName:%s ProjectID:%s}", c.Port, c.Env, c.CacheBucketName, c.ProjectID())
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

	githubApp, err := resolveGitHubApp()
	if err != nil {
		return nil, err
	}
	config.GitHubApp = githubApp
	if config.GitHubApp == nil && config.GitHubAuthToken == "" {
		return nil, fmt.Errorf("either GITHUB_APP_ID/GITHUB_APP_INSTALLATION_ID/GITHUB_APP_PRIVATE_KEY or GITHUB_AUTH_TOKEN is required")
	}

	config.CacheBucketName = os.Getenv("CACHE_STORAGE_BUCKET")
	config.googleCredentials = findGoogleCredentials()
	return &config, nil
}

// resolveGitHubApp reads the three GitHub App environment variables. It returns
// nil, nil when none of them are set, so the caller can fall back to
// GITHUB_AUTH_TOKEN for local development against an unregistered app. Setting
// only some of the three is treated as a misconfiguration rather than silently
// falling back, since that would hide a mistake.
func resolveGitHubApp() (*GitHubAppConfig, error) {
	appIDStr := os.Getenv("GITHUB_APP_ID")
	installationIDStr := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	privateKey := os.Getenv("GITHUB_APP_PRIVATE_KEY")

	set := 0
	for _, v := range []string{appIDStr, installationIDStr, privateKey} {
		if v != "" {
			set++
		}
	}
	if set == 0 {
		return nil, nil
	}
	if set != 3 {
		return nil, fmt.Errorf("GITHUB_APP_ID, GITHUB_APP_INSTALLATION_ID, and GITHUB_APP_PRIVATE_KEY must all be set together")
	}

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("GITHUB_APP_ID must be numeric: %w", err)
	}
	installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("GITHUB_APP_INSTALLATION_ID must be numeric: %w", err)
	}

	return &GitHubAppConfig{
		AppID:          appID,
		InstallationID: installationID,
		PrivateKey:     privateKey,
	}, nil
}

func findGoogleCredentials() *google.Credentials {
	cred, _ := google.FindDefaultCredentials(context.Background())
	return cred
}
