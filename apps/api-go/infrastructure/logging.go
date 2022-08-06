package infrastructure

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"golang.org/x/oauth2/google"
)

type LoggingClient struct {
	*logging.Client
}

func NewLoggingClient() *LoggingClient {
	cred, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(err)
	}
	if cred.ProjectID == "" {
		panic(fmt.Errorf("no project id found"))
	}
	c, err := logging.NewClient(context.Background(), cred.ProjectID)
	if err != nil {
		panic(err)
	}
	return &LoggingClient{c}
}
