package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoute(router *gin.Engine) {
	router.POST("/Login", controllers.Login())
	router.POST("/SignUp", controllers.SignUp())
	router.GET("/Refresh", controllers.Refresh())
	router.POST("/Logout", controllers.Logout())
}
