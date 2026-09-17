// Package github provides GitHub API client functionality
package github

import (
	"context"
	"fmt"
	"net/http"

	"contrib.rocks/apps/api/go/httptrace"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

// provider holds a single *github.Client. A *github.Client is safe for
// concurrent use, and the GitHub App transport (ghinstallation) caches its
// installation token internally, so building more than one client per
// process would cause redundant token issuance for no benefit.
type provider struct {
	client *github.Client
}

// NewProvider builds a provider authenticated as a GitHub App installation.
// The returned transport refreshes the installation token before it expires.
// It returns an error if privateKey is not a valid PEM-encoded RSA key.
func NewProvider(appID, installationID int64, privateKey string) (*provider, error) {
	tr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, []byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub App installation transport: %w", err)
	}
	return &provider{
		client: github.NewClient(&http.Client{Transport: httptrace.NewTransport(tr)}),
	}, nil
}

// NewTokenProvider builds a provider authenticated with a static personal
// access token. It is used as a local-development fallback when the GitHub
// App is not configured.
func NewTokenProvider(token string) *provider {
	oc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	oc.Transport = httptrace.NewTransport(oc.Transport)
	return &provider{
		client: github.NewClient(oc),
	}
}

func (f *provider) Get() *github.Client {
	return f.client
}
