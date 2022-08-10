package setting

import (
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/spf13/viper"
)

type Setting struct {
	V *viper.Viper
}

//数据库配置
type DbSettings struct {
	Host         string `form:"dbHost" binding:"required,max=50" json:"host"`
	Port         int    `form:"dbPort" binding:"required,gte=1,lte=65535" json:"port"`
	User         string `form:"dbUsername" binding:"required,max=50" json:"user"`
	Password     string `form:"dbPassword" binding:"required,max=50" json:"password"`
	Database     string `form:"dbName" binding:"required,max=30" json:"database"`
	Prefix       string `form:"dbTablePrefix" binding:"required,max=20" json:"prefix"`
	Charset      string `json:"charset"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type CacheSetting bigcache.Config

func NewSetting() *Setting {
	vp := viper.New()
	vp.AddConfigPath("configs/")
	vp.SetConfigName("configs")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		logger.NewLogger.Fatal().Msgf("read config error:{}", err)
	}
	return &Setting{V: vp}
}

//采用单个key并将其解组为struct
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.V.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
