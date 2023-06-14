package appcache

import (
	"context"
	"encoding/json"

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
func (cache *memoryCache) Save(c context.Context, name string, data []byte, contentType string) (model.FileHandle, error) {
	cache.fileCache[name] = &memoryFileHandle{name}
	return cache.fileCache[name], nil
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
	name string
}

func (f *memoryFileHandle) DownloadURL() string {
	return f.name
}
