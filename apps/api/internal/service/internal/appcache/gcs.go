package appcache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/sync/errgroup"
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
	return getFile(s.bucket, c, name)
}

func (s *gcsCache) GetJSON(c context.Context, name string, v any) error {
	o, err := getFile(s.bucket, c, name)
	if err != nil {
		return err
	}
	if o == nil {
		v = nil
		return nil
	}
	r := o.Reader()
	defer r.Close()
	return json.NewDecoder(r).Decode(&v)
}

func (s *gcsCache) Save(c context.Context, name string, data []byte, contentType string) error {
	return saveFile(s.bucket, c, name, data, contentType)
}

func (s *gcsCache) SaveJSON(c context.Context, name string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return saveFile(s.bucket, c, name, data, "application/json")
}

func getFile(bucket *storage.BucketHandle, c context.Context, name string) (model.FileHandle, error) {
	if bucket == nil {
		return nil, nil
	}
	ctx, span := tracing.Tracer().Start(c, "cache.Service.getFile")
	defer span.End()
	span.SetAttributes(attribute.String("cache.object.name", name))

	obj := bucket.Object(name)
	eg, ctx := errgroup.WithContext(ctx)
	file := &gcsFileHandle{}
	eg.Go(func() error {
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			return err
		}
		file.attrs = attrs
		return nil
	})
	eg.Go(func() error {
		or, err := obj.NewReader(ctx)
		if err != nil {
			return err
		}
		file.r = or
		return nil
	})
	err := eg.Wait()
	if err == storage.ErrObjectNotExist {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return file, nil
}

func saveFile(bucket *storage.BucketHandle, c context.Context, name string, data []byte, contentType string) error {
	if bucket == nil {
		return nil
	}
	ctx, span := tracing.Tracer().Start(c, "cache.Service.saveFile")
	defer span.End()
	span.SetAttributes(attribute.String("cache.object.name", name))

	w := bucket.Object(name).NewWriter(ctx)
	defer w.Close()
	w.ContentType = contentType
	_, err := w.Write(data)
	return err
}

var _ model.FileHandle = &gcsFileHandle{}

type gcsFileHandle struct {
	r     *storage.Reader
	attrs *storage.ObjectAttrs
}

func (h *gcsFileHandle) Reader() io.ReadCloser {
	return h.r
}
func (h *gcsFileHandle) Size() int64 {
	return h.attrs.Size
}
func (h *gcsFileHandle) ContentType() string {
	return h.attrs.ContentType
}
func (h *gcsFileHandle) ETag() string {
	return fmt.Sprintf("%x", h.attrs.MD5)
}
