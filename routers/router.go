package routers

import (
	"easyweb/routers/v1/sessions"
	"easyweb/routers/v1/users"
	"easyweb/utils/setting"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	//router.Use(gin.Logger())
	//router.Use(gin.Recovery())

	gin.SetMode(setting.Mode)

	user := router.Group("/v1/users")
	{
		// 注册
		user.POST("", users.SignUp)

		// 关注
		user.POST("/:id", users.Star)
		user.DELETE("/:id", users.UnStar)
		user.GET("/:id/stars", users.GetStars)
		user.GET("/:id/fans", users.GetFans)
	}

	session := router.Group("/v1/sessions")
	{
		// 登录退出
		session.POST("", sessions.SignIn)
		session.DELETE("", sessions.SignOut)
	}

	return router
}
