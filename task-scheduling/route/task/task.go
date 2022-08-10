package task

import (
	"encoding/json"
	"net/http"

	"github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/pkg"
	"github.com/gin-gonic/gin"
)

const (
	topic = "topic"
)

//创建topic
func TopicCreate(ctx *gin.Context) {
	topicInstance := new(models.Topic)
	if err := ctx.ShouldBindJSON(topicInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topicInstance.ID = pkg.GetID()
	topicByte, err := json.Marshal(topicInstance)
	//todo 存入缓存map[topicName:topicID]topic
	if err != nil {
		logger.NewLogger.Fatal().Msgf("to json err", err)
	}

}

func taskUpload(ctx *gin.Context) {

}
