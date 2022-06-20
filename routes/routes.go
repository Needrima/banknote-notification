package routes

import (
	"banknote-notification-service/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, controller *controllers.NotificationController) {
	router.POST("/notification", controller.SendNotification())
	router.GET("/notification/:reference/status", controller.GetNotificationStatus())
	router.GET("/notification/:reference", controller.GetNotificationByReference())
	router.GET("/notification/page/:page", controller.GetScheduledNotificationList())
}
