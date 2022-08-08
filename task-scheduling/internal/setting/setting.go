package setting

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Setting struct {
	V *viper.Viper
}

//数据库配置
type DbSettings struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	Prefix       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewSetting() *Setting {
	vp := viper.New()
	vp.AddConfigPath("configs/")
	vp.SetConfigName("configs")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatalf("read config error", err)
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
