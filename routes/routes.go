package routes

import (
	"banknote-notification-service/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *handler.NotificationHandler) {
	router.POST("/notification", handler.SendNotification())
	router.GET("/notification/:reference/status", handler.GetNotificationStatus())
	router.GET("/notification/:reference", handler.GetNotificationByReference())
	router.GET("/notification/page/:page", handler.GetScheduledNotificationList())
}
