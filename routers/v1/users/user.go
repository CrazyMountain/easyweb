package users

import (
	"easyweb/models"
	"easyweb/routers/v1/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	// 校验用户名
	if len(username) == 0 {
		errMsg := "Field username missing."
		common.OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	if ok, _ := models.IsUserExists(username); ok {
		errMsg := fmt.Sprintf("User %s already exists.", username)
		common.OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	// 校验密码
	password := c.PostForm("password")
	if len(password) == 0 {
		errMsg := "Field password missing."
		common.OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	// 用户入库
	if err := models.AddUser(username, password); nil != err {
		common.OperationFailed(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建成功
	description := fmt.Sprintf("User %s signed up successfully.", username)
	common.OperationSuccess(c, description, "")
}
