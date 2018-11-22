package sessions

import (
	"easyweb/models"
	"easyweb/routers/v1/common"
	"github.com/gin-gonic/gin"
)

func SignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if !common.ValidateUsernameAndPassword(username, password, c) {
		return
	}
	if !common.CheckPassword(username, password, c) {
		return
	}
	models.StartSession(username)
}

func SignOut(c *gin.Context) {

}
