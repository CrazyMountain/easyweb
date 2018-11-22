package common

import (
	"easyweb/models"
	"easyweb/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 校验参数字段
func ValidateFiled(field, value string, c *gin.Context) bool {
	if len(value) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   response.Msg[http.StatusBadRequest],
			"field":   field,
			"message": fmt.Sprintf("Field %s missing.", field),
		})
		return false
	}
	return true
}

// 校验用户名和密码
func ValidateUsernameAndPassword(username, password string, c *gin.Context) bool {
	if !ValidateFiled("username", username, c) {
		return false
	}
	if models.IsUserExists(username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   response.Msg[http.StatusBadRequest],
			"message": fmt.Sprintf("User %s already exists.", username),
		})
		return false
	}
	if !ValidateFiled("password", password, c) {
		return false
	}
	return true
}

// 检查密码是否正确
func CheckPassword(username, password string, c *gin.Context) bool {
	if models.IsPasswordCorrect(username, password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   response.Msg[http.StatusBadRequest],
			"message": "Password incorrect.",
		})
		return false
	}
	return true
}
