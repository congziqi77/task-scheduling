package routers

import "github.com/gin-gonic/gin"


//创建gin实例
func NewRouter() *gin.Engine {
	r := gin.Default()
	return r
}
