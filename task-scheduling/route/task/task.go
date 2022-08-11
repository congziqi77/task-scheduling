package task

import (
	"net/http"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/pkg"
	"github.com/gin-gonic/gin"
)

//创建topic
func TopicCreate(ctx *gin.Context) {
	topicInstance := new(models.Topic)
	if err := ctx.ShouldBindJSON(topicInstance); err != nil {
		logger.Debug().Str("err", err.Error()).Msg("bind err")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topicInstance.ID = pkg.GetID()
	err := topicInstance.SaveTopic2Cache()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
	})
}

//获取topic列表
func TopicList(ctx *gin.Context) {
	b, err := global.BigCache.Get(models.TopicKey)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("get topic from cache error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  string(b),
	})
}
