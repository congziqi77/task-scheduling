package task

import (
	"net/http"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/route/req"
	"github.com/gin-gonic/gin"
)

//创建topic
func TopicCreate(ctx *gin.Context) {
	topicInstance := new(models.Topic)
	if err := ctx.ShouldBindJSON(topicInstance); err != nil {
		logger.Error().Str("err", err.Error()).Msg("bind err")
		req.Error(ctx, err.Error())
		return
	}

	err := topicInstance.SaveTopic2Cache()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Success(ctx, "success")
}

//获取topic列表
func TopicList(ctx *gin.Context) {
	b, err := global.BigCache.Get(models.TopicKey)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("get topic from cache error")
		req.Error(ctx, err.Error())
		return
	}
	req.Success(ctx, string(b))
}

//创建task
func TaskCreate(ctx *gin.Context) {
	tasks := new(models.Tasks)
	if err := ctx.ShouldBindJSON(tasks); err != nil {
		logger.Error().Str("err", err.Error()).Msg("bind err")
		req.Error(ctx, err.Error())
		return
	}
	topicName := ctx.Query("topicName")
	err := tasks.TaskCreateServer(topicName)
	if err != nil {
		req.Error(ctx, err.Error())
		return
	}
	req.Success(ctx, "success")
}

//todo 增加查询依赖的接口 如果返回err或者出现没有值的情况那么重新执行同步接口
