package routes

import (
	"net/smtp"
	docs "walls-notification-service/docs"
	"walls-notification-service/internal/adapter/api"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
	errorhelper "walls-notification-service/internal/core/helper/error-helper"
	logger "walls-notification-service/internal/core/helper/log-helper"
	message "walls-notification-service/internal/core/helper/message-helper"
	"walls-notification-service/internal/core/middleware"
	"walls-notification-service/internal/core/services"
	ports "walls-notification-service/internal/port"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/twilio/twilio-go"
)

func SetupRouter(notificationRepository ports.NotificationRepository, redisClient *redis.Client, smtpClient *smtp.Client, smsClient *twilio.RestClient) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	notificationService := services.NewNotification(notificationRepository, redisClient, smtpClient, smsClient)

	handler := api.NewHTTPHandler(notificationService)

	logger.LogEvent("INFO", "Configuring Routes!")
	router.Use(middleware.LogRequest)

	corrs_config := cors.DefaultConfig()
	corrs_config.AllowAllOrigins = true

	router.Use(cors.New(corrs_config))
	//router.Use(middleware.SetHeaders)

	docs.SwaggerInfo.Description = "Walls Notification Service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = configuration.ServiceConfiguration.ServiceName

	router.POST("/api/notification", handler.CreateNotification)
	router.GET("/api/notification/user/:user-reference", handler.GetNotificationByUserReference)
	router.GET("/api/notification/device/:device-reference/:page", handler.GetNotificationByDeviceReference)
	router.GET("/api/notification/:reference", handler.GetNotificationByReference)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404,
			errorhelper.ErrorMessage(errorhelper.NoResourceError, message.NoResourceFound))
	})

	return router
}
