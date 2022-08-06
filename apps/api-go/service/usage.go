package service

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api-go/core"
	"contrib.rocks/apps/api-go/infrastructure"
	"contrib.rocks/libs/goutils/config"
	"contrib.rocks/libs/goutils/model"
)

type UsageService struct {
	env config.Environment
	l   *infrastructure.LoggingClient
}

func NewUsageService(i *core.Infrastructure) *UsageService {
	return &UsageService{i.Env, i.LoggingClient}
}

func (s *UsageService) Collect(ctx context.Context, r *model.RepositoryContributors, via string) error {
	logger := s.l.Logger("repository-usage")
	logger.Log(logging.Entry{
		Labels: map[string]string{
			"environment": string(s.env),
			"via":         via,
		},
		Payload: struct {
			Repository   string `json:"repository"`
			Owner        string `json:"owner"`
			RepoName     string `json:"repo"`
			Stargazers   int    `json:"stargazers"`
			Contributors int    `json:"contributors"`
			Timestamp    string `json:"timestamp"`
		}{
			Repository:   r.Repository.String(), // TODO remove
			Owner:        r.Repository.Owner,
			RepoName:     r.Repository.RepoName,
			Stargazers:   r.StargazersCount,
			Contributors: len(r.Contributors),
			Timestamp:    fmt.Sprint(time.Now().UnixMilli()), // TODO remove
		},
	})

	return nil
}
