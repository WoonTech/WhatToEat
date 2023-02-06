package routes

import (
	"what-to-eat/controllers"

	"github.com/gin-gonic/gin"
)

func PollRoute(router *gin.Engine) {
	router.POST("/Polls", controllers.CreatePoll())
	router.GET("/Polls/:id", controllers.GetPoll())
	router.GET("/Polls", controllers.GetAllPoll())
	router.DELETE("/Polls/:id", controllers.DeletePoll())
	router.PUT("/Polls/:id", controllers.UpdatePoll())
}
