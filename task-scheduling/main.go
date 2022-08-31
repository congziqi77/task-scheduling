package main

import (
	"net/http"
	"time"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/internal/routers"
	"github.com/congziqi77/task-scheduling/internal/setting"

	"github.com/gin-gonic/gin"
)

func init() {
	//初始化log
	logger.LogInit(true)

	if err := setupSetting(); err != nil {
		logger.Error().Str("err", err.Error()).Msg("init set error : %v")
		panic(err)
	}
	models.CacheImp = models.NewCache()
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

	if err = set.ReadSection("server", &global.ServerSetting); err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if err = set.ReadSection("database", &global.DbSetting); err != nil {
		return nil
	}
	return nil
}
