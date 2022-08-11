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
		logger.Error().Str("err", err.Error()).Msg("bind err")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
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
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  string(b),
	})
}

func TaskCreate(ctx *gin.Context) {
	tasks := new(models.Tasks)
	if err := ctx.ShouldBindJSON(tasks); err != nil {
		logger.Error().Str("err", err.Error()).Msg("bind err")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	topicName := ctx.Query("topicName")
	for _, task := range tasks.TaskList {
		task.ID = pkg.GetID()
		task.TopicName = topicName
	}
	maps, err := models.GetTopicMapFromCache()
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("get cache err")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	topic := maps[topicName]
	topic.Tasks = append(topic.Tasks, tasks.TaskList...)
	err = models.SetTopicMapToCache(maps)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("set cache error")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
	})
}
