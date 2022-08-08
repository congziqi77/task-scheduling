package install

import (
	"github.com/congziqi77/task-scheduling/internal/modules/app"
	"github.com/gin-gonic/gin"
)

type InstallForm struct {
	DBHost        string `binding:"Required;MaxSize(50)"`
	DBPort        int    `binding:"Required;Range(1,65535)"`
	DBUsername    string `binding:"Required;MaxSize(50)"`
	DBPassword    string `binding:"Required;MaxSize(50)"`
	DBName        string `binding:"Required;MaxSize(30)"`
	DBTablePrefix string `binding:"Required;MaxSize(20)"`
}

const (
	MAX_IDLE_Conns int = 10
	MAX_OPEN_CONNS int = 100
)

//绑定DBSetting
func DbBind(ctx *gin.Context, form InstallForm) {
	app.DbSetting.Host = form.DBHost
	app.DbSetting.Port = form.DBPort
	app.DbSetting.User = form.DBUsername
	app.DbSetting.Password = form.DBPassword
	app.DbSetting.Database = form.DBName
	app.DbSetting.Prefix = form.DBTablePrefix
	app.DbSetting.MaxIdleConns = MAX_IDLE_Conns
	app.DbSetting.MaxOpenConns = MAX_OPEN_CONNS
}
