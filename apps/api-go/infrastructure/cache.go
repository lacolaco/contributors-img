package infrastructure

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"contrib.rocks/libs/goutils/model"
)

var _ model.FileHandle = &cacheFileHandle{}

type cacheFileHandle struct {
	r *storage.Reader
}

func (h *cacheFileHandle) Reader() io.ReadCloser {
	return h.r
}
func (h *cacheFileHandle) Size() int64 {
	return h.r.Attrs.Size
}
func (h *cacheFileHandle) ContentType() string {
	return h.r.Attrs.ContentType
}
func (h *cacheFileHandle) ContentLength() int64 {
	return h.r.Attrs.Size
}

type CacheStorage struct {
	bucket *storage.BucketHandle
}

func NewCacheStorage() *CacheStorage {
	var cs CacheStorage
	bucketName := os.Getenv("CACHE_STORAGE_BUCKET")
	if bucketName == "" {
		return &cs
	}
	c, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	cs.bucket = c.Bucket(bucketName)
	return &cs
}

func (c *CacheStorage) Get(ctx context.Context, name string) (model.FileHandle, error) {
	return c.getFile(ctx, name)
}

func (c *CacheStorage) GetJSON(ctx context.Context, name string, v any) error {
	o, err := c.getFile(ctx, name)
	if err != nil {
		return err
	}
	if o == nil {
		v = nil
		return nil
	}
	return json.NewDecoder(o.Reader()).Decode(&v)
}

func (c *CacheStorage) Save(ctx context.Context, name string, data []byte, contentType string) error {
	return c.saveFile(ctx, name, data, contentType)
}

func (c *CacheStorage) SaveJSON(ctx context.Context, name string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.saveFile(ctx, name, data, "application/json")
}

func (c *CacheStorage) getFile(ctx context.Context, name string) (model.FileHandle, error) {
	if c.bucket == nil {
		return nil, nil
	}
	or, err := c.bucket.Object(name).NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &cacheFileHandle{or}, nil
}

func (c *CacheStorage) saveFile(ctx context.Context, name string, data []byte, contentType string) error {
	if c.bucket == nil {
		return nil
	}
	w := c.bucket.Object(name).NewWriter(ctx)
	defer w.Close()
	w.ContentType = contentType
	_, err := w.Write(data)
	return err
}
