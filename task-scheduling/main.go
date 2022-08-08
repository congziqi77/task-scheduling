package main

import (
	"log"
	"net/http"
	"time"

	"github.com/congziqi77/task-scheduling/internal/modules/app"
	"github.com/congziqi77/task-scheduling/internal/routers"
	"github.com/congziqi77/task-scheduling/internal/setting"
	"github.com/gin-gonic/gin"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init error : %v", err)
	}
	//todo 测试数据库连接
}

func main() {
	gin.SetMode(app.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + app.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    app.ServerSetting.ReadTimeout,
		WriteTimeout:   app.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func setupSetting() error {
	set := setting.NewSetting()
	err := set.ReadSection("Server", &app.ServerSetting)
	if err != nil {
		return err
	}
	app.ServerSetting.ReadTimeout *= time.Second
	app.ServerSetting.WriteTimeout *= time.Second
	return nil
}
