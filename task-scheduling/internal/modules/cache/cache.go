package cache

import (
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
)

func CacheInit() {
	config := global.CacheSetting
	var cache *bigcache.BigCache
	var initErr error
	logger.NewLogger.Info().Msgf("cache config is: %v", config)
	if config == nil {
		cache, initErr = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	} else {
		cache, initErr = bigcache.NewBigCache(bigcache.Config(*config))
	}
	if initErr != nil {
		logger.NewLogger.Fatal().Msgf("bigCache init error", initErr)
	}
	global.BigCache = cache
}

//todo 公共创建缓存的方法 和 json -》map-》struct

