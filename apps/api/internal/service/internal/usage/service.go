package usage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/env"
	"contrib.rocks/libs/goutils/model"
)

type Service struct {
	env           env.Environment
	loggingClient *logging.Client
}

func New(cfg *config.Config, l *logging.Client) *Service {
	return &Service{cfg.Env, l}
}

func (s *Service) CollectUsage(c context.Context, r *model.RepositoryContributors, via string) error {
	_, span := tracing.DefaultTracer.Start(c, "usage.Service.CollectUsage")
	defer span.End()

	logger := s.loggingClient.Logger("repository-usage")
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
