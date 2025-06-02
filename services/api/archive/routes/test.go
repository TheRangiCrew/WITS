package routes

import "github.com/gin-gonic/gin"

func RegisterTestRoutes(router *gin.Engine) {
	// Test route to check if the server is running
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
		})
	})

	// Additional test routes can be added here
}
