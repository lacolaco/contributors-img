package appcache

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"

	"contrib.rocks/libs/go/model"
)

var _ AppCache = &memoryCache{}

// memoryCache is a simple in-memory cache. (for local debug or testing)
type memoryCache struct {
	fileCache map[string]model.FileHandle
	jsonCache map[string][]byte
}

func newMemoryCache() *memoryCache {
	return &memoryCache{
		fileCache: make(map[string]model.FileHandle),
		jsonCache: make(map[string][]byte),
	}
}

func (cache *memoryCache) Get(c context.Context, name string) (model.FileHandle, error) {
	return cache.fileCache[name], nil
}
func (cache *memoryCache) GetJSON(c context.Context, name string, v any) error {
	cached := cache.jsonCache[name]
	if cached == nil {
		return nil
	}
	return json.Unmarshal(cached, &v)
}
func (cache *memoryCache) Save(c context.Context, name string, data []byte, contentType string) error {
	cache.fileCache[name] = &memoryFileHandle{data, contentType}
	return nil
}
func (cache *memoryCache) SaveJSON(c context.Context, name string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	cache.jsonCache[name] = data
	return nil
}

var _ model.FileHandle = &memoryFileHandle{}

type memoryFileHandle struct {
	data        []byte
	contentType string
}

func (f *memoryFileHandle) ContentType() string {
	return f.contentType
}

func (f *memoryFileHandle) ETag() string {
	return fmt.Sprintf("%x", md5.Sum(f.data))
}

func (f *memoryFileHandle) Reader() io.ReadCloser {
	return io.NopCloser(bytes.NewBuffer(f.data))
}

func (f *memoryFileHandle) Size() int64 {
	return int64(len(f.data))
}
