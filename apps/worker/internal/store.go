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

func SaveFeaturedRepositories(ctx context.Context, appEnv env.Environment, usages []*RepositoryUsage, updatedAt time.Time) error {
	fs := apiclient.NewFirestoreClient()
	items := make([]FeaturedRepository, len(usages))
	for i, repository := range usages {
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
