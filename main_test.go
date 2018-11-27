package main

import (
	"easyweb/routers/v1"
	"encoding/json"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type failedResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type successResponse struct {
	Data        interface{}
	Description string `json:"description"`
	Message     string `json:"message"`
}

var router *gin.Engine

func init() {
	router = gin.Default()

	user := router.Group("/v1/users")
	{
		// 注册
		user.POST("", v1.SignUp)
	}

	follow := user.Group("/:username/follows")
	{
		// 关注
		follow.POST("", v1.Follow)
		follow.DELETE("", v1.UnFollow)
		follow.GET("/:flag", v1.GetFollows)
	}

	session := router.Group("/v1/sessions")
	{
		// 登录退出
		session.POST("", v1.SignIn)
		session.DELETE("", v1.SignOut)
	}
}

func TestSignUp(t *testing.T) {

	body := url.Values{"username": []string{"neo"}}

	request := httptest.NewRequest("POST", "/v1/users", strings.NewReader(body.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()

	/*Patch(models.IsUserExists, func(_ string) (bool, error) {
		return false, nil
	})
	defer UnpatchAll()*/

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	Convey("Test sign up.", t, func() {

		Convey("Test user already exists.", func() {
			bytes, _ := ioutil.ReadAll(result.Body)
			var f failedResponse
			json.Unmarshal(bytes, &f)

			So(result.StatusCode, ShouldEqual, http.StatusBadRequest)
			So(f.Error, ShouldEqual, "User neo already exists.")
		})

		Convey("", func() {

		})

	})

}
