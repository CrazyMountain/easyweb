package routers

import (
	"easyweb/routers/v1"
	"easyweb/utils/setting"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	gin.SetMode(setting.Mode)

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

	return router
}
