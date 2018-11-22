package users

import (
	"easyweb/models"
	"easyweb/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Field username missing.",
			"error":   response.Msg[http.StatusBadRequest],
			"field":   "username",
		})
		return
	}

	password := c.PostForm("password")
	if len(password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Field password missing.",
			"error":   response.Msg[http.StatusBadRequest],
			"field":   "password",
		})
		return
	}

	models.AddUser(username, password)

	c.JSON(http.StatusOK, gin.H{
		"message": response.Msg[http.StatusOK],
	})
}
