package cache

import (
	"context"
	"encoding/json"
	"io"

	"cloud.google.com/go/storage"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/model"
)

type Service struct {
	bucket *storage.BucketHandle
}

func New(cfg *config.Config, storageClient *storage.Client) *Service {
	var cs Service
	if cfg.CacheBucketName != "" {
		cs.bucket = storageClient.Bucket(cfg.CacheBucketName)
	}
	return &cs
}

func (c *Service) Get(ctx context.Context, name string) (model.FileHandle, error) {
	return getFile(c.bucket, ctx, name)
}

func (c *Service) GetJSON(ctx context.Context, name string, v any) error {
	o, err := getFile(c.bucket, ctx, name)
	if err != nil {
		return err
	}
	if o == nil {
		v = nil
		return nil
	}
	defer o.Reader().Close()
	b, err := io.ReadAll(o.Reader())
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}

func (c *Service) Save(ctx context.Context, name string, data []byte, contentType string) error {
	return saveFile(c.bucket, ctx, name, data, contentType)
}

func (c *Service) SaveJSON(ctx context.Context, name string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return saveFile(c.bucket, ctx, name, data, "application/json")
}

func getFile(bucket *storage.BucketHandle, ctx context.Context, name string) (model.FileHandle, error) {
	if bucket == nil {
		return nil, nil
	}
	or, err := bucket.Object(name).NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &cacheFileHandle{or}, nil
}

func saveFile(bucket *storage.BucketHandle, ctx context.Context, name string, data []byte, contentType string) error {
	if bucket == nil {
		return nil
	}
	w := bucket.Object(name).NewWriter(ctx)
	defer w.Close()
	w.ContentType = contentType
	_, err := w.Write(data)
	return err
}
