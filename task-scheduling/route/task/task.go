package task

import (
	"github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/route/req"
	"github.com/gin-gonic/gin"
)

// 创建task
func TaskCreate(ctx *gin.Context) {
	tasks := new(models.Tasks)
	if err := ctx.ShouldBindJSON(tasks); err != nil {
		logger.Error().Str("err", err.Error()).Msg("bind err")
		req.Error(ctx, err.Error())
		return
	}
	topicID := ctx.Query("topicID")
	topicName := ctx.Query("topicName")

	if err := tasks.TaskCreateServer(topicName, topicID); err != nil {
		logger.Error().Str("err", err.Error()).Msg("")
		req.Error(ctx, err.Error())
		return
	}
	req.Success(ctx, "success")
}
