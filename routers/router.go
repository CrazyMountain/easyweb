package routers

import (
	"easyweb/routers/v1/sessions"
	"easyweb/routers/v1/users"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	user := router.Group("/v1/user")
	{
		user.POST("", users.SignUp)

		// star
		user.POST("/:id", users.Star)
		user.DELETE("/:id", users.UnStar)
		user.GET("/:id/stars", users.GetStars)
		user.GET("/:id/fans", users.GetFans)
	}

	session := router.Group("/v1/session")
	{
		// sign in
		session.POST("", sessions.SignIn)
		session.POST("/:id", sessions.SignOut)
	}

	return router
}
