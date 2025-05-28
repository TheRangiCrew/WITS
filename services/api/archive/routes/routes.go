package routes

import "github.com/gin-gonic/gin"

func Setup() *gin.Engine {
	router := gin.Default()

	// Register routes
	RegisterTestRoutes(router)
	RegisterVTECRoutes(router)

	return router
}
