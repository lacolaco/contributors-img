package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

func SaveFeaturedRepositories(ctx context.Context, repositories []FeaturedRepository, updatedAt time.Time) error {
	env := GetEnvironment(ctx)
	c, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		return err
	}
	items := make([]RepositoryUsage, len(repositories))
	for i, repository := range repositories {
		items[i] = repository.Usage
	}
	_, err = c.Collection(env).Doc("featured_repositories_2").Set(ctx, struct {
		Items     []RepositoryUsage `firestore:"items"`
		UpdatedAt time.Time         `firestore:"updatedAt"`
	}{
		items, updatedAt,
	})
	return err
}
