package models

import (
	"runtime/debug"

	"github.com/congziqi77/task-scheduling/internal/modules/inter"

	"github.com/coocood/freecache"
)

var CacheImp inter.ICache

type LocalCache struct {
	Cache *freecache.Cache
}

func cacheInit() *freecache.Cache {
	cacheSize := 100 * 1024 * 1024
	debug.SetGCPercent(20)
	cache := freecache.NewCache(cacheSize)
	return cache
}

func NewCache() *LocalCache {
	return &LocalCache{
		Cache: cacheInit(),
	}
}

func (cache *LocalCache) SetCache(key []byte, val []byte, expireSeconds int) error {
	return cache.Cache.Set(key, val, expireSeconds)
}

func (cache *LocalCache) GetCache(key []byte) ([]byte, error) {
	b, err := cache.Cache.Get(key)
	return b, err
}

func (cache *LocalCache) ClearCache() {
	cache.Cache.Clear()
}
