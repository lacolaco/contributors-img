package apiclient

import (
	"context"
	"net/http"

	"contrib.rocks/libs/goutils/httptrace"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient(token string) *github.Client {
	oc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	oc.Transport = httptrace.NewTransport(http.DefaultTransport)
	c := github.NewClient(oc)
	return c
}
