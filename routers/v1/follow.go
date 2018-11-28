package v1

import (
	"easyweb/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Follow(c *gin.Context) {
	fan := c.Param("username")
	star := c.Param("followed")

	// 校验用户是否存在
	if ok, _ := models.IsUserExists(fan); !ok {
		errMsg := fmt.Sprintf("User %s does not exists.", fan)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	if ok, _ := models.IsUserExists(star); !ok {
		errMsg := fmt.Sprintf("User %s does not exists.", star)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	// 校验是否已经关注
	if ok, _ := models.IsFollowExists(fan, star); ok {
		errMsg := fmt.Sprintf("User %s has already followed user %s.", fan, star)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	if err := models.AddFollow(fan, star); nil != err {
		OperationFailed(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 操作成功
	description := fmt.Sprintf("User %s follows user %s.", fan, star)
	OperationSuccess(c, description, "")
}

func UnFollow(c *gin.Context) {
	fan := c.Param("username")
	star := c.Param("followed")

	// 校验用户是否存在
	if ok, _ := models.IsUserExists(fan); !ok {
		errMsg := fmt.Sprintf("User %s does not exists.", fan)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	if ok, _ := models.IsUserExists(star); !ok {
		errMsg := fmt.Sprintf("User %s does not exists.", star)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	// 校验是否已经关注
	if ok, _ := models.IsFollowExists(fan, star); !ok {
		errMsg := fmt.Sprintf("User %s did not follow user %s.", fan, star)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	if err := models.DeleteFollow(fan, star); nil != err {
		OperationFailed(c, http.StatusInternalServerError, err.Error())
		return
	}

	description := fmt.Sprintf("User %s unfollows user %s.", fan, star)
	OperationSuccess(c, description, "")
}

func GetFollows(c *gin.Context) {
	username := c.Param("username")
	flag := c.Param("flag")

	// 校验用户是否存在
	if ok, _ := models.IsUserExists(username); !ok {
		errMsg := fmt.Sprintf("User %s does not exists.", username)
		OperationFailed(c, http.StatusBadRequest, errMsg)
		return
	}

	switch flag {
	case "0":
		fans, err := models.GetFollows(username)
		if nil != err {
			OperationFailed(c, http.StatusInternalServerError, err.Error())
			return
		}
		description := fmt.Sprintf("User %s' followers.", username)
		OperationSuccess(c, description, fans)
	case "1":
		stars, err := models.GetFollowed(username)
		if nil != err {
			OperationFailed(c, http.StatusInternalServerError, err.Error())
			return
		}
		description := fmt.Sprintf("Users followed by user %s.", username)
		OperationSuccess(c, description, stars)
	default:
		errMsg := fmt.Sprintf("Illegal parameter: %s.", flag)
		OperationFailed(c, http.StatusBadRequest, errMsg)
	}
}
