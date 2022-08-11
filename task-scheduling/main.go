package main

import (
	"net/http"
	"time"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/cache"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/internal/routers"
	"github.com/congziqi77/task-scheduling/internal/setting"

	"github.com/gin-gonic/gin"
)

func init() {
	//初始化log
	logger.Loginit(true)
	err := setupSetting()
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("init set error : %v")
		panic(err)
	}
	err = cache.CacheInit()
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("init cache error : %v")
		panic(err)
	}
}

func main() {

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func setupSetting() error {
	set, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = set.ReadSection("server", &global.ServerSetting)
	if err != nil {
		logger.Debug().Str("err", err.Error()).Msg("read yaml error")
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	err = set.ReadSection("cache", &global.CacheSetting)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("read yaml error :")
		return err
	}
	return nil
}
