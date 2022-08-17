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
		logger.Error().Str("err", err.Error()).Msg("save topic 2 cache error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Success(ctx, "success")
}

//获取topic列表
func TopicList(ctx *gin.Context) {
	b, err := global.FreeCache.Get([]byte(models.TopicKey))
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("get topic from cache error")
		req.Error(ctx, err.Error())
		return
	}
	req.Success(ctx, string(b))
}

//TODO 增加查询依赖的接口 如果返回err或者出现没有值的情况那么重新执行同步接口
func GetTopo(ctx *gin.Context) {
	topicID := ctx.Query("topicID")
	topicName := ctx.Query("topicName")
	s, err := models.GetTopicTopo(topicName, topicID)
	if err != nil {
		logger.Error().Str("get topo cache error:", err.Error()).Msg("error")
		req.Error(ctx, "get topo cache error")
	}
	ctx.JSON(http.StatusOK, gin.H{"topoRes": s})
}


func Run(ctx *gin.Context) {
	topicName := ctx.Query("topicName")
	topicID := ctx.Query("topicID")

	models.Run(topicName,topicID)
}