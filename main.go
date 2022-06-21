package main

import (
	"banknote-notification-service/database"
	"banknote-notification-service/handler"
	"banknote-notification-service/middlewares"
	"banknote-notification-service/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	logger     *log.Logger
	collection *mongo.Collection
)

func init() {
	logger = log.New(os.Stdout, "====>  ", log.LstdFlags)

	collection = database.InitCollection("notifications")
}

func main() {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	controller := handler.New(logger, collection)

	routes.SetupRoutes(router, controller)

	router.Run(":8080")
}
