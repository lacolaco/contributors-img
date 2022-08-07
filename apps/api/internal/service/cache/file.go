package cache

import (
	"io"

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
