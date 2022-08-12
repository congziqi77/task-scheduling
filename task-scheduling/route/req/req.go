package req

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  message,
	})
}

func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusBadRequest,
		"msg":    message,
	})
}
