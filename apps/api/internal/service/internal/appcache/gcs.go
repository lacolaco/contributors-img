package appcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/internal/tracing"
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

// isNotExistErr reports whether err means the object is absent, which this
// cache treats as a miss rather than a failure. storage wraps
// ErrObjectNotExist rather than returning it bare, so equality does not work.
func isNotExistErr(err error) bool {
	return errors.Is(err, storage.ErrObjectNotExist)
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

	// Using Background context to prevent cancellation during fetching operations
	bgCtx := context.Background()
	eg, _ := errgroup.WithContext(bgCtx)

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

	if err := eg.Wait(); err != nil {
		if file.reader != nil {
			file.reader.Close()
		}
		if isNotExistErr(err) {
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
	w.ContentType = contentType
	if _, err := w.Write(data); err != nil {
		w.Close()
		return err
	}
	// Write only fills the writer's internal pipe. Whether the upload succeeded is
	// reported by Close alone, so its error is the one that matters here.
	return w.Close()
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
