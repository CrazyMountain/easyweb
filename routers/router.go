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

	gin.SetMode(setting.Mode)

	user := router.Group("/v1/user")
	{
		// 注册
		user.POST("", users.SignUp)

		// 关注
		user.POST("/:id", users.Star)
		user.DELETE("/:id", users.UnStar)
		user.GET("/:id/stars", users.GetStars)
		user.GET("/:id/fans", users.GetFans)
	}

	session := router.Group("/v1/session")
	{
		// 登录退出
		session.POST("", sessions.SignIn)
		session.POST("/:id", sessions.SignOut)
	}

	return router
}
