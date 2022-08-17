package cache

import (
	"runtime/debug"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/coocood/freecache"
)

func CacheInit() {
	cacheSize := 100 * 1024 * 1024
	debug.SetGCPercent(20)
	cache := freecache.NewCache(cacheSize)
	global.FreeCache = *cache
}
