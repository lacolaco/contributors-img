package app

import (
	"context"
	"time"

	"contrib.rocks/libs/go/apiclient"
	"contrib.rocks/libs/go/env"
)

type FeaturedRepository struct {
	Repository   string `firestore:"repository"`
	Stargazers   int64  `firestore:"stargazers"`
	Contributors int64  `firestore:"contributors"`
}

func SaveFeaturedRepositories(ctx context.Context, appEnv env.Environment, usageRows []*RepositoryUsageRow, updatedAt time.Time) error {
	fs := apiclient.NewFirestoreClient()
	items := make([]FeaturedRepository, len(usageRows))
	for i, repository := range usageRows {
		items[i] = FeaturedRepository{
			Repository:   repository.Repository,
			Stargazers:   repository.Stargazers,
			Contributors: repository.Contributors,
		}
	}
	_, err := fs.Collection(string(appEnv)).Doc("featured_repositories").Set(ctx, struct {
		Items     []FeaturedRepository `firestore:"items"`
		UpdatedAt time.Time            `firestore:"updatedAt"`
	}{
		items, updatedAt,
	})
	return err
}

type UsageStats struct {
	Owners       int64 `firestore:"owners"`
	Repositories int64 `firestore:"repositories"`
}

func SaveUsageStats(ctx context.Context, appEnv env.Environment, usageRow *UsageStatsRow, updatedAt time.Time) error {
	fs := apiclient.NewFirestoreClient()

	_, err := fs.Collection(string(appEnv)).Doc("usage_stats").Set(ctx, struct {
		Owners       int64     `firestore:"owners"`
		Repositories int64     `firestore:"repositories"`
		UpdatedAt    time.Time `firestore:"updatedAt"`
	}{
		usageRow.Owners, usageRow.Repositories, updatedAt,
	})
	return err
}
