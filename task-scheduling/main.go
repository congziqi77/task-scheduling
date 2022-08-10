package main

import (
	"log"
	"net/http"
	"time"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/internal/routers"
	"github.com/congziqi77/task-scheduling/internal/setting"

	"github.com/gin-gonic/gin"
)

func init() {
	//初始化log
	logger.NewLogger = logger.Loginit(true)
	err := setupSetting()
	if err != nil {
		log.Fatalf("init error : %v", err)
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
	set := setting.NewSetting()
	err := set.ReadSection("server", &global.ServerSetting)
	if err != nil {
		logger.NewLogger.Fatal().Msgf("read yaml error :", err)
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	err = set.ReadSection("cache", &global.CacheSetting)
	if err != nil {
		logger.NewLogger.Fatal().Msgf("read yaml error :", err)
		return err
	}
	return nil
}


