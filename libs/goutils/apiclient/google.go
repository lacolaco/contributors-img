package apiclient

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
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

func NewLoggingClient() *logging.Client {
	cred, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(err)
	}
	if cred.ProjectID == "" {
		panic(fmt.Errorf("project id is not found"))
	}
	c, err := logging.NewClient(context.Background(), cred.ProjectID)
	if err != nil {
		panic(err)
	}
	return c
}
