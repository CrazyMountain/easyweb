package users

import (
	"easyweb/models"
	"easyweb/routers/v1/common"
	"easyweb/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 校验用户名和密码
	if !common.ValidateUsernameAndPassword(username, password, c) {
		return
	}

	// 用户入库
	models.AddUser(username, password)

	c.JSON(http.StatusOK, gin.H{
		"message": response.Msg[http.StatusOK],
	})
}
