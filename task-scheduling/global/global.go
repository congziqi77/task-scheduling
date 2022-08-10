package global

import (
	"github.com/allegro/bigcache/v3"
	"github.com/congziqi77/task-scheduling/internal/setting"
)

var (
	DbSetting     *setting.DbSettings
	ServerSetting *setting.ServerSettingS
	CacheSetting  *setting.CacheSetting
)

var (
	BigCache *bigcache.BigCache
)
