package routers

import (
	"github.com/congziqi77/task-scheduling/route/install"
	"github.com/congziqi77/task-scheduling/route/task"
	"github.com/gin-gonic/gin"
)

//创建gin实例
func NewRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/task/conn")
	{
		apiV1.POST("/conndb", install.DbBind)
	}
	apiV2 := r.Group("/task")
	{
		apiV2.POST("/topicCreate", task.TopicCreate)
		apiV2.GET("/topicList", task.TopicList)
	}
	return r
}
