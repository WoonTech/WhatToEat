package main

import (
	"what-to-eat/configuration"
	"what-to-eat/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "Hello"})
	})
	configuration.ConnectDB()
	routes.PollRoute(router)
	routes.ResRoute(router)
	routes.UserRoute(router)

	router.Run("localhost:6000")
}
