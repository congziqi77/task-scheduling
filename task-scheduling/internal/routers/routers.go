package routers

import (
	"github.com/congziqi77/task-scheduling/route/install"
	"github.com/congziqi77/task-scheduling/route/task"
	"github.com/gin-gonic/gin"
)

// 创建gin实例
func NewRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/task/conn")
	{
		apiV1.POST("/connDB", install.DbBind)
	}
	apiV2 := r.Group("/task/topic")
	{
		apiV2.POST("/create", task.TopicCreate)
		apiV2.GET("/list", task.TopicList)
		apiV2.GET("/getTopo", task.GetTopo)
		apiV2.GET("/run", task.Run)
	}
	apiV3 := r.Group("/task")
	{
		apiV3.POST("/create", task.TaskCreate)
	}
	return r
}
