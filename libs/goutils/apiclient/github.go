package apiclient

import (
	"context"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient(token string) *github.Client {
	c := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))
	return c
}
