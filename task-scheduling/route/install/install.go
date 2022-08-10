package install

import (
	"net/http"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	MAX_IDLE_Conns int = 10
	MAX_OPEN_CONNS int = 100
)

//绑定DBSetting并开启DB
func DbBind(ctx *gin.Context) {
	if err := ctx.ShouldBind(&global.DbSetting); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	global.DbSetting.MaxIdleConns = MAX_IDLE_Conns
	global.DbSetting.MaxOpenConns = MAX_OPEN_CONNS
	models.DB = models.NewDBEngine()
	ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})

}
