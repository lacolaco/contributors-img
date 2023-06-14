package appcache

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/storage"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/go/model"
	"go.opentelemetry.io/otel/attribute"
)

var _ AppCache = &gcsCache{}

type gcsCache struct {
	bucket *storage.BucketHandle
}

func newGCSCache(storageClient *storage.Client, bucketName string) *gcsCache {
	return &gcsCache{
		bucket: storageClient.Bucket(bucketName),
	}
}

func (s *gcsCache) Get(c context.Context, name string) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "appcache.Get")
	defer span.End()

	return getFile(s.bucket, ctx, name)
}

func (s *gcsCache) Save(c context.Context, name string, data []byte, contentType string) (model.FileHandle, error) {
	ctx, span := tracing.Tracer().Start(c, "appcache.Save")
	defer span.End()
	return saveFile(s.bucket, ctx, name, data, contentType)
}

func (s *gcsCache) GetJSON(c context.Context, name string, v any) error {
	ctx, span := tracing.Tracer().Start(c, "appcache.GetJSON")
	defer span.End()
	return getJSON(s.bucket, ctx, name, v)
}

func (s *gcsCache) SaveJSON(c context.Context, name string, v any) error {
	ctx, span := tracing.Tracer().Start(c, "appcache.SaveJSON")
	defer span.End()
	return saveJSON(s.bucket, ctx, name, v)
}

func getFile(bucket *storage.BucketHandle, c context.Context, name string) (model.FileHandle, error) {
	if bucket == nil {
		return nil, nil
	}
	ctx, span := tracing.Tracer().Start(c, "appcache.getFile")
	defer span.End()

	span.SetAttributes(attribute.String("cache.object.name", name))

	// Context from errgroup will be canceled whether error is returned or not. But the GCS client need a context to run RPCs.
	// So we use a separate context.
	_, err := bucket.Object(name).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, nil
		}
		return nil, err
	}
	url, err := getSignedURL(bucket, ctx, name)
	if err != nil {
		return nil, err
	}

	return &gcsFileHandle{url}, nil
}

func saveFile(bucket *storage.BucketHandle, c context.Context, name string, data []byte, contentType string) (model.FileHandle, error) {
	if bucket == nil {
		return nil, nil
	}
	ctx, span := tracing.Tracer().Start(c, "appcache.saveFile")
	defer span.End()
	span.SetAttributes(attribute.String("cache.object.name", name))

	obj := bucket.Object(name)
	w := obj.NewWriter(ctx)
	defer w.Close()
	w.ContentType = contentType
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	url, err := getSignedURL(bucket, ctx, name)
	if err != nil {
		return nil, err
	}
	return &gcsFileHandle{url}, nil
}

func getSignedURL(bucket *storage.BucketHandle, c context.Context, name string) (string, error) {
	_, span := tracing.Tracer().Start(c, "appcache.getSignedURL")
	defer span.End()

	url, err := bucket.SignedURL(name, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return "", err
	}
	return url, nil
}

func getJSON(bucket *storage.BucketHandle, c context.Context, name string, v any) error {
	if bucket == nil {
		return nil
	}
	ctx, span := tracing.Tracer().Start(c, "appcache.getJSON")
	defer span.End()

	span.SetAttributes(attribute.String("cache.object.name", name))

	obj := bucket.Object(name)
	r, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil
		}
		return err
	}
	return json.NewDecoder(r).Decode(&v)
}

func saveJSON(bucket *storage.BucketHandle, c context.Context, name string, v any) error {
	if bucket == nil {
		return nil
	}
	ctx, span := tracing.Tracer().Start(c, "appcache.saveJSON")
	defer span.End()
	span.SetAttributes(attribute.String("cache.object.name", name))

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w := bucket.Object(name).NewWriter(ctx)
	defer w.Close()
	w.ContentType = "application/json"
	if _, err = w.Write(data); err != nil {
		return err
	}
	return nil
}

var _ model.FileHandle = &gcsFileHandle{}

type gcsFileHandle struct {
	signedURL string
}

func (h *gcsFileHandle) DownloadURL() string {
	return h.signedURL
}
