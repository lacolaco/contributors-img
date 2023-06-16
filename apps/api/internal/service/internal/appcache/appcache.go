package appcache

import (
	"context"

	"cloud.google.com/go/storage"
	"contrib.rocks/libs/go/model"
)

type AppCache interface {
	Get(c context.Context, name string) (model.FileHandle, error)
	GetJSON(c context.Context, name string, v any) error
	Save(c context.Context, name string, data []byte, contentType string) (model.FileHandle, error)
	SaveJSON(c context.Context, name string, v any) error
}

func NewGCSCache(storageClient *storage.Client, bucketName string) AppCache {
	return newGCSCache(storageClient, bucketName)
}

func NewMemoryCache() AppCache {
	return newMemoryCache()
}
