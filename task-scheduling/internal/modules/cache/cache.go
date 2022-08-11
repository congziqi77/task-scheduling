package cache

import (
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/congziqi77/task-scheduling/global"
)

func CacheInit() error {
	config := global.CacheSetting
	var cache *bigcache.BigCache
	var err error
	if config == nil {
		cache, err = bigcache.NewBigCache(bigcache.DefaultConfig(1 * time.Hour))
	} else {
		cache, err = bigcache.NewBigCache(bigcache.Config(*config))
	}
	if err != nil {
		return err
	}
	global.BigCache = cache
	return nil
}
