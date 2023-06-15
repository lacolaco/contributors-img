package appcache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/go/model"
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
	ctx, span := tracing.Tracer().Start(c, "appcache.Get")
	defer span.End()

	return getFile(s.bucket, ctx, name)
}

func (s *gcsCache) GetJSON(c context.Context, name string, v any) error {
	ctx, span := tracing.Tracer().Start(c, "appcache.GetJSON")
	defer span.End()

	o, err := getFile(s.bucket, ctx, name)
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
	ctx, span := tracing.Tracer().Start(c, "appcache.Save")
	defer span.End()
	return saveFile(s.bucket, ctx, name, data, contentType)
}

func (s *gcsCache) SaveJSON(c context.Context, name string, v any) error {
	ctx, span := tracing.Tracer().Start(c, "appcache.SaveJSON")
	defer span.End()

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return saveFile(s.bucket, ctx, name, data, "application/json")
}

func getFile(bucket *storage.BucketHandle, c context.Context, name string) (model.FileHandle, error) {
	if bucket == nil {
		return nil, nil
	}
	ctx, span := tracing.Tracer().Start(c, "appcache.getFile")
	defer span.End()

	span.SetAttributes(attribute.String("cache.object.name", name))

	obj := bucket.Object(name)
	file := &gcsFileHandle{}
	// Context from errgroup will be canceled whether error is returned or not. But the GCS client need a context to run RPCs.
	// So we use a separate context.
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			return err
		}
		file.attrs = attrs
		return nil
	})
	eg.Go(func() error {
		r, err := obj.NewReader(ctx)
		if err != nil {
			return err
		}
		file.reader = r
		return nil
	})
	err := eg.Wait()
	if err != nil {
		if file.reader != nil {
			file.reader.Close()
		}
		if err == storage.ErrObjectNotExist {
			return nil, nil
		}
		return nil, err
	}
	return file, nil
}

func saveFile(bucket *storage.BucketHandle, c context.Context, name string, data []byte, contentType string) error {
	if bucket == nil {
		return nil
	}
	ctx, span := tracing.Tracer().Start(c, "appcache.saveFile")
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
	reader *storage.Reader
	attrs  *storage.ObjectAttrs
}

func (h *gcsFileHandle) Reader() io.ReadCloser {
	return h.reader
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
