package global

import (
	"github.com/congziqi77/task-scheduling/internal/setting"
	"github.com/coocood/freecache"
	"gorm.io/gorm"
)

var (
	DbSetting     *setting.DbSettings
	ServerSetting *setting.ServerSettingS
)

var (
	FreeCache freecache.Cache
	DB        *gorm.DB
)

const (
	TopicTopoSuffix       = "Topo"
	ISStartGetFromResChan = "ISStartGetFromResChan"
)
