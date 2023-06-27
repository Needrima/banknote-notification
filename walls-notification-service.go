package main

import (
	"context"
	"walls-notification-service/internal/adapter/events/subscriber"
	extensions "walls-notification-service/internal/adapter/extensions"
	mongoRepository "walls-notification-service/internal/adapter/repository/mongodb"

	"fmt"
	"walls-notification-service/internal/adapter/routes"
	channel "walls-notification-service/internal/core/domain/event/channel"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
	logger "walls-notification-service/internal/core/helper/log-helper"
	message "walls-notification-service/internal/core/helper/message-helper"
)

func main() {
	//Initialize request Log
	logger.InitializeLog()

	//Start DB Connection
	mongoRepo := extensions.StartDatabase("mongodb")
	logger.LogEvent("INFO", "MongoDB Connected and Initialized!")

	//start redis
	logger.LogEvent("INFO", message.StartingRedis)
	redisClient := extensions.StartEventBus("redis")

	// start twilio sms client
	smsClient := extensions.StartTwilioConnection("twilio")
	logger.LogEvent("INFO", "Initialized twilio for sms!")

	// start sendgrid email client
	emailClient := extensions.StartSendGridConnection("sendgrid")
	logger.LogEvent("INFO", "Initialized sendgrid for email!")

	//Set up routes
	router := routes.SetupRouter(mongoRepo.(mongoRepository.MongoRepositories).Notification, redisClient, smsClient, emailClient)

	config := configuration.ServiceConfiguration

	go func() {
		logger.LogEvent("INFO", message.StartingServer)
		err := router.Run(":" + config.ServicePort)
		//api.SetConfiguration
		if err != nil {
			fmt.Println(err)
			logger.LogEvent("ERROR", "Error Starting Server : "+err.Error())
		}
	}()

	// Initialize the event subscriber
	eventSubscriber := subscriber.NewEventSubscriber(redisClient)
	ctx := context.Background()
	// Run the subscription code in a Goroutine
	go func() {
		eventSubscriber.SubscribeToSignUpCreatedEvent(ctx, channel.OtpCreatedEvent)
	}()

	select {}
}
