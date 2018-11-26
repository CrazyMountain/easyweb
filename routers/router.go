package routers

import (
	"easyweb/routers/v1/sessions"
	"easyweb/routers/v1/users"
	"easyweb/utils/setting"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	gin.LoggerWithWriter()
	gin.SetMode(setting.Mode)

	user := router.Group("/v1/users")
	{
		// 注册
		user.POST("", users.SignUp)
	}

	follow := user.Group("/:username/follows")
	{
		// 关注
		follow.POST("", users.Follow)
		follow.DELETE("", users.UnFollow)
		follow.GET("/:flag", users.GetFollows)
	}

	session := router.Group("/v1/sessions")
	{
		// 登录退出
		session.POST("", sessions.SignIn)
		session.DELETE("", sessions.SignOut)
	}

	return router
}
