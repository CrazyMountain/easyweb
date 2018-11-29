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

// some const for test
const (
	testUsername  = "neo"
	testPassword  = "123"
	testSessionID = "bfv0msjjp201ka19i6mg"
	testStar      = "mary"
	randomFlag    = "5"
)

// api targets for test
var targets = map[string]string{
	"users":    "/v1/users",
	"sessions": "/v1/sessions",
	"follows":  "/v1/users/" + testUsername + "/follows/" + testStar,
	"fans":     "/v1/users/" + testUsername + "/follows/0",
	"stars":    "/v1/users/" + testUsername + "/follows/1",
	"defaults": "/v1/users/" + testUsername + "/follows/" + randomFlag,
}

type failedResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type successResponse struct {
	Data        interface{} `json:"data"`
	Description string      `json:"description"`
	Message     string      `json:"message"`
}

type testParameters struct {
	bodies       []url.Values
	descriptions []string
	statuses     []int
	errors       []string
}

var signUp, signIn, signOut, followT, unFollowT, follows, defaults testParameters

var router *gin.Engine

func init() {
	router = gin.New()

	user := router.Group("/v1/users")
	{
		// sign up
		user.POST("", v1.SignUp)
	}

	follow := user.Group("/:username/follows")
	{
		// follow
		follow.POST("/:followed", v1.Follow)
		follow.DELETE("/:followed", v1.UnFollow)
		follow.GET("/:flag", v1.GetFollows)
	}

	session := router.Group("/v1/sessions")
	{
		// sign int and sign out
		session.POST("", v1.SignIn)
		session.DELETE("", v1.SignOut)
	}

	// initial test parameters

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
		bodies: []url.Values{},
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

	followT = testParameters{
		bodies: []url.Values{},
		descriptions: []string{
			"Test user fan does not exist",
			"Test user star does not exist",
			"Test already followed",
			"Test failed to operate database",
			"Test success",
		},
		statuses: []int{
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusOK,
		},
		errors: []string{
			fmt.Sprintf("User %s does not exists.", testUsername),
			fmt.Sprintf("User %s does not exists.", testStar),
			fmt.Sprintf("User %s has already followed user %s.", testUsername, testStar),
			"database error",
		},
	}

	unFollowT = testParameters{
		bodies: []url.Values{},
		descriptions: []string{
			"Test user fan does not exist",
			"Test user star does not exist",
			"Test not followed",
			"Test failed to operate database",
			"Test success",
		},
		statuses: []int{
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusOK,
		},
		errors: []string{
			fmt.Sprintf("User %s does not exists.", testUsername),
			fmt.Sprintf("User %s does not exists.", testStar),
			fmt.Sprintf("User %s did not follow user %s.", testUsername, testStar),
			"database error",
		},
	}

	follows = testParameters{
		bodies: []url.Values{},
		descriptions: []string{
			"Test user does not exist",
			"Test failed to operate database",
			"Test success",
		},
		statuses: []int{
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusOK,
		},
		errors: []string{
			fmt.Sprintf("User %s does not exists.", testUsername),
			"database error",
		},
	}

	defaults = testParameters{
		bodies: []url.Values{},
		descriptions: []string{
			"Test user does not exist",
			"Test failed",
		},
		statuses: []int{
			http.StatusBadRequest,
			http.StatusBadRequest,
		},
		errors: []string{
			fmt.Sprintf("User %s does not exists.", testUsername),
			fmt.Sprintf("Illegal parameter: %s.", randomFlag),
		},
	}
}

func TestSignUp(t *testing.T) {
	defer UnpatchAll()

	Convey("Test sign up.", t, func() {
		for i, desc := range signUp.descriptions {

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

			request := httptest.NewRequest(http.MethodPost, targets["users"], strings.NewReader(signUp.bodies[i].Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 6 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s signed up successfully.", testUsername))
					So(result.StatusCode, ShouldEqual, signUp.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, signUp.statuses[i])
				So(r.Error, ShouldEqual, signUp.errors[i])
			})
		}
	})
}

func TestSignIn(t *testing.T) {
	defer UnpatchAll()

	Convey("Test sign in.", t, func() {
		for i, desc := range signIn.descriptions {

			if i == 0 {
				Patch(models.IsSignIn, func(_ string) (bool, error) {
					return true, nil
				})
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

			request := httptest.NewRequest(http.MethodPost, targets["sessions"], strings.NewReader(signIn.bodies[i].Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if i == 0 {
				c := http.Cookie{
					Name:  "session_id",
					Value: testSessionID,
				}
				request.AddCookie(&c)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 8 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s signed in with session %s.", testUsername, testSessionID))
					So(result.StatusCode, ShouldEqual, signIn.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, signIn.statuses[i])
				So(r.Error, ShouldEqual, signIn.errors[i])
			})
		}
	})
}

func TestSignOut(t *testing.T) {
	defer UnpatchAll()

	Convey("Test sign out.", t, func() {
		for i, desc := range signOut.descriptions {

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

			request := httptest.NewRequest(http.MethodDelete, targets["sessions"], nil)

			if i > 0 {
				c := http.Cookie{
					Name:  "session_id",
					Value: testSessionID,
				}
				request.AddCookie(&c)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 2 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("Session %s deleted.", testSessionID))
					So(result.StatusCode, ShouldEqual, signOut.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, signOut.statuses[i])
				So(r.Error, ShouldEqual, signOut.errors[i])
			})
		}
	})
}

func TestFollow(t *testing.T) {
	defer UnpatchAll()

	Convey("Test follow.", t, func() {
		for i, desc := range followT.descriptions {

			if i == 0 {
				Patch(models.IsUserExists, func(username string) (bool, error) {
					if username == "neo" {
						return false, nil
					}
					return true, nil
				})
			}
			if i == 1 {
				Patch(models.IsUserExists, func(username string) (bool, error) {
					if username == "neo" {
						return true, nil
					}
					return false, nil
				})
			}
			if i == 2 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
				Patch(models.IsFollowExists, func(_, _ string) (bool, error) {
					return true, nil
				})
			}
			if i == 3 {
				Patch(models.IsFollowExists, func(_, _ string) (bool, error) {
					return false, nil
				})
				Patch(models.AddFollow, func(_, _ string) error {
					return fmt.Errorf("database error")
				})
			}
			if i == 4 {
				Patch(models.AddFollow, func(_, _ string) error {
					return nil
				})
			}

			request := httptest.NewRequest(http.MethodPost, targets["follows"], nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 4 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s follows user %s.", testUsername, testStar))
					So(result.StatusCode, ShouldEqual, followT.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, followT.statuses[i])
				So(r.Error, ShouldEqual, followT.errors[i])
			})
		}
	})
}

func TestUnFollow(t *testing.T) {
	defer UnpatchAll()

	Convey("Test unFollow.", t, func() {
		for i, desc := range unFollowT.descriptions {

			if i == 0 {
				Patch(models.IsUserExists, func(username string) (bool, error) {
					if username == "neo" {
						return false, nil
					}
					return true, nil
				})
			}
			if i == 1 {
				Patch(models.IsUserExists, func(username string) (bool, error) {
					if username == "neo" {
						return true, nil
					}
					return false, nil
				})
			}
			if i == 2 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
				Patch(models.IsFollowExists, func(_, _ string) (bool, error) {
					return false, nil
				})
			}
			if i == 3 {
				Patch(models.IsFollowExists, func(_, _ string) (bool, error) {
					return true, nil
				})
				Patch(models.DeleteFollow, func(_, _ string) error {
					return fmt.Errorf("database error")
				})
			}
			if i == 4 {
				Patch(models.DeleteFollow, func(_, _ string) error {
					return nil
				})
			}

			request := httptest.NewRequest(http.MethodDelete, targets["follows"], nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 4 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(r.Data, ShouldEqual, "")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s unfollows user %s.", testUsername, testStar))
					So(result.StatusCode, ShouldEqual, unFollowT.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, unFollowT.statuses[i])
				So(r.Error, ShouldEqual, unFollowT.errors[i])
			})
		}
	})
}

func TestGetFollows(t *testing.T) {
	defer UnpatchAll()

	Convey("Test get fans.", t, func() {
		for i, desc := range follows.descriptions {

			if i == 0 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return false, nil
				})
			}
			if i == 1 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
				Patch(models.GetFollows, func(_ string) ([]string, error) {
					return []string{}, fmt.Errorf("database error")
				})
			}
			if i == 2 {
				Patch(models.GetFollows, func(_ string) ([]string, error) {
					return []string{"mary", "mike"}, nil
				})
			}

			request := httptest.NewRequest(http.MethodGet, targets["fans"], nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 2 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(fmt.Sprint(r.Data), ShouldEqual, "[mary mike]")
					So(r.Description, ShouldEqual, fmt.Sprintf("User %s' followers.", testUsername))
					So(result.StatusCode, ShouldEqual, follows.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, follows.statuses[i])
				So(r.Error, ShouldEqual, follows.errors[i])
			})
		}
	})

	Convey("Test get stars.", t, func() {
		for i, desc := range follows.descriptions {

			if i == 0 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return false, nil
				})
			}
			if i == 1 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
				Patch(models.GetFollowed, func(_ string) ([]string, error) {
					return []string{}, fmt.Errorf("database error")
				})
			}
			if i == 2 {
				Patch(models.GetFollowed, func(_ string) ([]string, error) {
					return []string{"mary", "mike"}, nil
				})
			}

			request := httptest.NewRequest(http.MethodGet, targets["stars"], nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			if i == 2 {
				var r successResponse
				json.Unmarshal(bytes, &r)

				Convey(desc, func() {
					So(fmt.Sprint(r.Data), ShouldEqual, "[mary mike]")
					So(r.Description, ShouldEqual, fmt.Sprintf("Users followed by user %s.", testUsername))
					So(result.StatusCode, ShouldEqual, follows.statuses[i])
				})
				continue
			}

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, follows.statuses[i])
				So(r.Error, ShouldEqual, follows.errors[i])
			})
		}
	})

	Convey("Test default case.", t, func() {
		for i, desc := range defaults.descriptions {

			if i == 0 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return false, nil
				})
			}
			if i == 1 {
				Patch(models.IsUserExists, func(_ string) (bool, error) {
					return true, nil
				})
			}

			request := httptest.NewRequest(http.MethodGet, targets["defaults"], nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			result := recorder.Result()
			bytes, _ := ioutil.ReadAll(result.Body)

			var r failedResponse
			json.Unmarshal(bytes, &r)

			Convey(desc, func() {
				So(result.StatusCode, ShouldEqual, defaults.statuses[i])
				So(r.Error, ShouldEqual, defaults.errors[i])
			})
		}
	})
}
