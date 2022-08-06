package infrastructure

import (
	"context"
	"os"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	*github.Client
}

func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		panic("GITHUB_AUTH_TOKEN is not set")
	}
	c := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))
	return &GitHubClient{c}
}
