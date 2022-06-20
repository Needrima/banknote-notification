package routes

import (
	"banknote-notification-service/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/notification", controllers.SendNotification())
	router.GET("/notification/:reference/status", controllers.GetNotificationStatus())
	router.GET("/notification/:reference", controllers.GetNotificationByReference())
	router.GET("/notification/page/:page", controllers.GetScheduledNotificationList())
}
