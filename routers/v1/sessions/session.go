package sessions

import (
	"easyweb/models"
	"easyweb/routers/v1/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignIn(c *gin.Context) {
	// 是否已经登录
	id, err := c.Cookie("session_id")
	if nil == err {
		if ok, _ := models.IsSignIn(id); ok {
			errMsg := "Already signed in."
			common.OperationFailed(c, http.StatusInternalServerError, errMsg)
			return
		}
		fmt.Printf("Session %s does not exists or expired.\r\n", id)
	}

	// 校验用户名
	username := c.PostForm("username")
	if len(username) == 0 {
		errMsg := "Field username missing."
		common.OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	if ok, _ := models.IsUserExists(username); !ok {
		errMsg := fmt.Sprintf("User %s does not exists.", username)
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

	if ok, _ := models.IsPasswordCorrect(username, password); !ok {
		errMsg := "Password incorrect."
		common.OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	// 开启新会话
	newId, err := models.StartSession(username)
	if nil != err {
		common.OperationFailed(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("session_id", newId, 0, "/", "localhost", false, true)
	description := fmt.Sprintf("User %s signed in with session %s.", username, newId)
	common.OperationSuccess(c, description)
}

func SignOut(c *gin.Context) {
	id, err := c.Cookie("session_id")
	if nil != err {
		common.OperationFailed(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 退出登录
	if err = models.EndSession(id); nil != err {
		common.OperationFailed(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 操作成功
	description := fmt.Sprintf("Session %s deleted.", id)
	common.OperationSuccess(c, description)
}
