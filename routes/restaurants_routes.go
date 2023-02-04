package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func ResRoute(router *gin.Engine) {
	router.POST("/Restaurants", controllers.CreateRes())
	router.GET("/Restaurants/:id", controllers.GetRes())
	router.GET("/Restaurants", controllers.GetAllRes())
	router.DELETE("/Restaurants/:id", controllers.DeleteRes())
	router.PUT("/Restaurants/:id", controllers.UpdateRes())
}
