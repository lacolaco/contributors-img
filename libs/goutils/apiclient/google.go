package apiclient

import (
	"context"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

func NewBigQueryClient() *bigquery.Client {
	c, err := bigquery.NewClient(context.Background(), bigquery.DetectProjectID)
	if err != nil {
		panic(err)
	}
	return c
}

func NewStorageClient() *storage.Client {
	c, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	return c
}

func NewLoggingClient(projectID string) *logging.Client {
	c, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		panic(err)
	}
	return c
}

func NewFirestoreClient() *firestore.Client {
	c, err := firestore.NewClient(context.Background(), firestore.DetectProjectID)
	if err != nil {
		panic(err)
	}
	return c
}
