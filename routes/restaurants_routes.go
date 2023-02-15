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
	router.POST("/Menus", controllers.CreateMenu())
	router.GET("/Menus/:id", controllers.GetMenu())
	router.DELETE("/Menus/:id", controllers.DeleteMenu())
	router.PUT("/Menus/:id", controllers.UpdateMenu())
}
