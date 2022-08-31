package global

import (
	"github.com/congziqi77/task-scheduling/internal/setting"
	"gorm.io/gorm"
)

var (
	DbSetting     *setting.DbSettings
	ServerSetting *setting.ServerSettingS
)

var (
	DB *gorm.DB
)

const (
	TopicTopoSuffix       = "Topo"
	ISStartGetFromResChan = "ISStartGetFromResChan"
)
