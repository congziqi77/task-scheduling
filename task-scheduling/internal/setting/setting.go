package setting

import (
	"time"

	"github.com/spf13/viper"
)

type Setting struct {
	V *viper.Viper
}

// 数据库配置
type DbSettings struct {
	Host         string `binding:"required,max=50" json:"host"`
	Port         int    `binding:"required,gte=1,lte=65535" json:"port"`
	User         string `binding:"required,max=50" json:"user"`
	Password     string `binding:"required,max=50" json:"password"`
	Database     string `binding:"required,max=30" json:"database"`
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

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.AddConfigPath("/Users/congziqi/Documents/goWork/task-scheduling/task-scheduling/configs")
	vp.SetConfigName("configs")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	return &Setting{V: vp}, nil
}

// 采用单个key并将其解组为struct
func (s *Setting) ReadSection(k string, v interface{}) error {
	if err := s.V.UnmarshalKey(k, v); err != nil {
		return err
	}
	return nil
}
