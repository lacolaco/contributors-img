package app

import (
	"context"
	"time"

	"contrib.rocks/libs/goutils/apiclient"
	"contrib.rocks/libs/goutils/env"
)

func SaveFeaturedRepositories(ctx context.Context, appEnv env.Environment, repositories []FeaturedRepository, updatedAt time.Time) error {
	fs := apiclient.NewFirestoreClient()
	items := make([]RepositoryUsage, len(repositories))
	for i, repository := range repositories {
		items[i] = repository.Usage
	}
	_, err := fs.Collection(string(appEnv)).Doc("featured_repositories").Set(ctx, struct {
		Items     []RepositoryUsage `firestore:"items"`
		UpdatedAt time.Time         `firestore:"updatedAt"`
	}{
		items, updatedAt,
	})
	return err
}
