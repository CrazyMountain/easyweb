package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Success = "Operation success."
	Failed  = "Operation failed."
)

func OperationSuccess(c *gin.Context, description string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"description": description,
		"message":     Success,
		"data":        data,
	})
}

func OperationFailed(c *gin.Context, status int, errMsg string) {
	c.JSON(status, gin.H{
		"error":   errMsg,
		"message": Failed,
	})
}
