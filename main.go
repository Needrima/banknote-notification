package main

import (
	"banknote-notification-service/middlewares"
	"banknote-notification-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	routes.SetupRoutes(router)

	router.Run(":8080")
}
