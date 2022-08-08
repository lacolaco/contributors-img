package cache

import (
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"contrib.rocks/libs/goutils/model"
)

var _ model.FileHandle = &cacheFileHandle{}

type cacheFileHandle struct {
	r     *storage.Reader
	attrs *storage.ObjectAttrs
}

func (h *cacheFileHandle) Reader() io.ReadCloser {
	return h.r
}
func (h *cacheFileHandle) Size() int64 {
	return h.attrs.Size
}
func (h *cacheFileHandle) ContentType() string {
	return h.attrs.ContentType
}
func (h *cacheFileHandle) ETag() string {
	return fmt.Sprintf("%x", h.attrs.MD5)
}
