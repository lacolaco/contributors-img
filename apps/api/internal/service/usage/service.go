package usage

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api/internal/logger"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"go.uber.org/zap"
)

const (
	logGroupID = "repository-usage"
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
	log = log.With(logger.LogGroup(logGroupID), logger.Label("via", via))
	log.Info(fmt.Sprintf("repository-usage: %s", r.Repository.String()),
		zap.String("owner", r.Repository.Owner),
		zap.String("repository", r.Repository.String()),
		zap.Int("stargazers", r.StargazersCount),
		zap.Int("contributors", len(r.Contributors)),
	)
	return nil
}
