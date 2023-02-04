package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/Restaurants", controllers.CreateUser())
	router.GET("/Restaurants/:id", controllers.GetUser())
	router.GET("/Restaurants", controllers.GetAllUser())
	router.DELETE("/Restaurants/:id", controllers.DeleteUser())
	router.PUT("/Restaurants/:id", controllers.UpdateUser())
}
