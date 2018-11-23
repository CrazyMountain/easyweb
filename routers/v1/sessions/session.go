package sessions

import (
	"easyweb/models"
	"easyweb/routers/v1/common"
	"easyweb/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignIn(c *gin.Context) {
	username := c.PostForm("username")
	if !common.ValidateFiled("username", username, c) {
		return
	}
	if !models.IsUserExists(username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   response.Msg[http.StatusBadRequest],
			"message": fmt.Sprintf("User %s does not exists.", username),
		})
		return
	}

	password := c.PostForm("password")
	if !common.ValidateFiled("password", password, c) {
		return
	}
	if !common.CheckPassword(username, password, c) {
		return
	}

	id := models.StartSession(username)
	c.SetCookie("session_id", id, 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": response.Msg[http.StatusOK],
	})
}

func SignOut(c *gin.Context) {
	id, err := c.Cookie("session_id")
	if nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Cookie parse error.",
			"message": response.Msg[http.StatusInternalServerError],
		})
		return
	}
	if !models.IsSignIn(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Session does not exists or expired.",
			"message": response.Msg[http.StatusBadRequest],
		})
		return
	}

	models.StopSession(id)
	c.JSON(http.StatusOK, gin.H{
		"message": response.Msg[http.StatusOK],
	})
}
