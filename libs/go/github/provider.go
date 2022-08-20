package github

import (
	"context"
	"sync"

	"contrib.rocks/libs/go/httptrace"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type provider struct {
	pool sync.Pool
}

func NewProvider(token string) *provider {
	return &provider{
		pool: sync.Pool{
			New: func() any {
				oc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: token},
				))
				oc.Transport = httptrace.NewTransport(oc.Transport)
				return github.NewClient(oc)
			},
		},
	}
}

func (f *provider) Get() *github.Client {
	client := f.pool.Get().(*github.Client)
	f.pool.Put(client)
	return client
}
