package users

import (
	"easyweb/models"
	"easyweb/routers/v1/common"
	"easyweb/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	// 校验用户名和密码
	if !common.ValidateFiled("username", username, c) {
		return
	}
	if models.IsUserExists(username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   response.Msg[http.StatusBadRequest],
			"message": fmt.Sprintf("User %s already exists.", username),
		})
		return
	}

	password := c.PostForm("password")
	if !common.ValidateFiled("password", password, c) {
		return
	}

	// 用户入库
	models.AddUser(username, password)

	c.JSON(http.StatusOK, gin.H{
		"message": response.Msg[http.StatusOK],
	})
}
