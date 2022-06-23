package routes

import (
	"bnt/bnt-notification-service/internal/adapter/api"
	"bnt/bnt-notification-service/internal/core/helper"
	"bnt/bnt-notification-service/internal/core/middleware"
	"bnt/bnt-notification-service/internal/core/services"
	ports "bnt/bnt-notification-service/internal/port"
	"github.com/gin-gonic/gin"
)

func SetupRouter(notificationRepository ports.NotificationRepository) *gin.Engine {
	router := gin.Default()
	notificationService := services.NewNotification(notificationRepository)

	handler := api.NewHTTPHandler(notificationService)

	helper.LogEvent("INFO", "Configuring Routes!")
	router.Use(middleware.LogRequest)

	//router.Use(middleware.SetHeaders)

	router.Group("/notification")
	{
		router.POST("/notification", handler.CreateNotification)
		router.GET("/notification/:reference/status", handler.GetNotificationStatus)
		router.GET("/notification/page/:page", handler.GetNotificationList)
		router.GET("/notification/:reference", handler.GetNotificationByRef)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404,
			helper.ErrorMessage(helper.NoResourceError, helper.NoResourceFound))
	})
	return router
}
