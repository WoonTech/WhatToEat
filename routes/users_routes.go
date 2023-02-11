package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/Users", controllers.CreateUser())
	router.GET("/Users/:id", controllers.GetUser())
	router.GET("/Users", controllers.GetAllUser())
	router.DELETE("/Users/:id", controllers.DeleteUser())
	router.PUT("/Users/:id", controllers.UpdateUser())
}
