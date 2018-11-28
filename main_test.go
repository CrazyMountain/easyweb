package main

import (
	. "bou.ke/monkey"
	"easyweb/models"
	"easyweb/routers/v1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	testUsername  = "neo"
	testPassword  = "123"
	testSessionID = "bfv0msjjp201ka19i6mg"
)

var targets = map[string]string{
	"users":    "/v1/users",
	"follows":  "/v1/users/:username/follows",
	"sessions": "/v1/sessions",
}

type failedResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type successResponse struct {
	Data        interface{}
	Description string `json:"description"`
	Message     string `json:"message"`
}

type testParameters struct {
	bodies       []url.Values
	descriptions []string
	statuses     []int
	errors       []string
}

var signUp, signIn, signOut testParameters

var router *gin.Engine

func init() {
	router = gin.New()

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

	// 初始化测试参数

	signUp = testParameters{
		bodies: []url.Values{
			// Test field username missing.
			{},
			{
				"username": []string{""},
			},
			// Test username already exists.
			{
				"username": []string{testUsername},
			},
			// Test field password missing.
			{
				"username": []string{testUsername},
			},
			{
				"username": []string{testUsername},
				"password": []string{""},
			},
			// Test write to database.
			{
				"username": []string{testUsername},
				"password": []string{testPassword},
			},
			// Test success
			{
				"username": []string{testUsername},
				"password": []string{testPassword},
			},
		},
		descriptions: []string{
			"Test field username missing(non-field)",
			"Test field username missing(empty string)",
			"Test username already exists",
			"Test field password missing(non-field)",
			"Test field password missing(empty string)",
			"Test failed to write to database",
			"Test sign up successfully",
		},
		statuses: []int{
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusOK,
		},
		errors: []string{
			"Field username missing.",
			"Field username missing.",
			fmt.Sprintf("User %s already exists.", testUsername),
			"Field password missing.",
			"Field password missing.",
			"Database error...",
		},
	}

	signIn = testParameters{
		bodies: []url.Values{
			// Test session already exists.
			{},
			// Test field username missing.
			{},
			{
				"username": []string{""},
			},
			// Test username does not exists.
			{
				"username": []string{testUsername},
			},
			// Test field password missing.
			{
				"username": []string{testUsername},
			},
			{
				"username": []string{testUsername},
				"password": []string{""},
			},
			// Test password incorrect.
			{
				"username": []string{testUsername},
				"password": []string{testPassword},
			},
			// Test failed to write to database.
			{
				"username": []string{testUsername},
				"password": []string{testPassword},
			},
			// Test success.
			{
				"username": []string{testUsername},
				"password": []string{testPassword},
			},
		},
		descriptions: []string{
			"Test session already exists",
			"Test field username missing(non-field)",
			"Test field username missing(empty string)",
			"Test username does not exists",
			"Test field password missing(non-field)",
			"Test field password missing(empty string)",
			"Test password incorrect",
			"Test failed to write to database",
			"Test sign in successfully",
		},
		statuses: []int{
			http.StatusInternalServerError,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusOK,
		},
		errors: []string{
			"Already signed in.",
			"Field username missing.",
			"Field username missing.",
			fmt.Sprintf("User %s does not exists.", testUsername),
			"Field password missing.",
			"Field password missing.",
			"Password incorrect.",
			"Database error...",
		},
	}

	signOut = testParameters{
		bodies: []url.Values{
			// Test session does not exists.
			{},
			// Test failed to operate database.
			{},
			// Test success.
			{},
		},
		descriptions: []string{
			"Test session does not exists",
			"Test failed to operate database",
			"Test sign out successfully",
		},
		statuses: []int{
			http.StatusInternalServerError,
			http.StatusInternalServerError,
			http.StatusOK,
		},
		errors: []string{
			"http: named cookie not present",
			"database error",
		},
	}
}

func TestSignUp(t *testing.T) {
	defer UnpatchAll()
	Convey("Test sign up.", t, func() {

		for i, body := range signUp.bodies {
			if i == 2 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
			}
			if i == 3 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return false, nil
				})
			}
			if i == 5 {
				Patch(models.AddUser, func(_, _ string) error {
					return fmt.Errorf(signUp.errors[i])
				})
			}
			if i == 6 {
				Patch(models.AddUser, func(_, _ string) error {
					return nil
				})
			}

			request := httptest.NewRequest(http.MethodPost, targets["users"], strings.NewReader(body.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 6 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(signUp.descriptions[i], func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s signed up successfully.", testUsername))
					So(result.StatusCode, ShouldEqual, signUp.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(signUp.descriptions[i], func() {
				So(result.StatusCode, ShouldEqual, signUp.statuses[i])
				So(r.Error, ShouldEqual, signUp.errors[i])
			})
		}
	})
}

func TestSignIn(t *testing.T) {
	defer UnpatchAll()
	Convey("Test sign in.", t, func() {
		for i, body := range signIn.bodies {
			request := httptest.NewRequest(http.MethodPost, targets["sessions"], strings.NewReader(body.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if i == 0 {
				Patch(models.IsSignIn, func(_ string) (bool, error) {
					return true, nil
				})

				c := http.Cookie{
					Name:  "session_id",
					Value: testSessionID,
				}
				request.AddCookie(&c)
			}
			if i == 3 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return false, nil
				})
			}
			if i == 4 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
			}
			if i == 6 {
				Patch(models.IsPasswordCorrect, func(_, _ string) (bool, error) {
					return false, nil
				})
			}
			if i == 7 {
				Patch(models.IsPasswordCorrect, func(_, _ string) (bool, error) {
					return true, nil
				})
				Patch(models.StartSession, func(_ string) (string, error) {
					return testSessionID, fmt.Errorf(signIn.errors[i])
				})
			}
			if i == 8 {
				Patch(models.StartSession, func(_ string) (string, error) {
					return testSessionID, nil
				})
			}

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 8 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(signIn.descriptions[i], func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s signed in with session %s.", testUsername, testSessionID))
					So(result.StatusCode, ShouldEqual, signIn.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(signIn.descriptions[i], func() {
				So(result.StatusCode, ShouldEqual, signIn.statuses[i])
				So(r.Error, ShouldEqual, signIn.errors[i])
			})
		}
	})
}

func TestSignOut(t *testing.T) {
	defer UnpatchAll()
	Convey("Test sign out.", t, func() {
		for i, body := range signOut.bodies {

			request := httptest.NewRequest(http.MethodDelete, targets["sessions"], strings.NewReader(body.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if i > 0 {
				c := http.Cookie{
					Name:  "session_id",
					Value: testSessionID,
				}
				request.AddCookie(&c)
			}
			if i == 1 {
				Patch(models.EndSession, func(_ string) error {
					return fmt.Errorf("database error")
				})
			}
			if i == 2 {
				Patch(models.EndSession, func(_ string) error {
					return nil
				})
			}

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 2 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(signOut.descriptions[i], func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("Session %s deleted.", testSessionID))
					So(result.StatusCode, ShouldEqual, signOut.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(signOut.descriptions[i], func() {
				So(result.StatusCode, ShouldEqual, signOut.statuses[i])
				So(r.Error, ShouldEqual, signOut.errors[i])
			})
		}
	})
}

func TestFollow(t *testing.T) {

}

func TestUnFollow(t *testing.T) {

}

func TestGetFollows(t *testing.T) {

}
