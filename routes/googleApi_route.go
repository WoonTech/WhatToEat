package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func GoogleApiRoute(router *gin.Engine) {
	router.GET("/GoogleRestaurants", controllers.GetGoogleRes())
}
