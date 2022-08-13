package usage

import (
	"context"
	"time"

	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"go.uber.org/zap"
)

type Service interface {
	CollectUsage(c context.Context, r *model.RepositoryContributors, via string) error
}

func New() Service {
	return &serviceImpl{}
}

var _ Service = &serviceImpl{}

type serviceImpl struct {
}

func (s *serviceImpl) CollectUsage(c context.Context, r *model.RepositoryContributors, via string) error {
	ctx, span := tracing.Tracer().Start(c, "usage.Service.CollectUsage")
	defer span.End()

	log := logger.LoggerFromContext(ctx)
	log = log.With(
		logger.LogGroupLabel("repository-usage"),
		logger.Label("via", via),
	)
	log.Info(r.Repository.String(),
		zap.String("repository", r.Repository.String()),
		zap.String("owner", r.Repository.Owner),
		zap.String("repo", r.Repository.RepoName),
		zap.Int("stargazers", r.StargazersCount),
		zap.Int("contributors", len(r.Contributors)),
		zap.Time("timestamp", time.Now()),
	)
	return nil
}
