package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func ResRoute(router *gin.Engine) {
	router.POST("/Restaurants", controllers.CreateRes())
}
