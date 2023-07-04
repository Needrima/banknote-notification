package services

import (
	"context"
	"net/smtp"

	"time"
	publisher "walls-notification-service/internal/adapter/event/publisher"
	sender "walls-notification-service/internal/adapter/sender"
	"walls-notification-service/internal/core/domain/dto"
	event "walls-notification-service/internal/core/domain/event/eto"
	"walls-notification-service/internal/core/domain/mapper"
	"walls-notification-service/internal/core/domain/shared"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
	eto "walls-notification-service/internal/core/helper/event-helper/eto"
	logger "walls-notification-service/internal/core/helper/log-helper"

	ports "walls-notification-service/internal/port"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twilio/twilio-go"
)

var NotificationService = &notificationService{}

type notificationService struct {
	notificationRepository ports.NotificationRepository
	redisClient            *redis.Client
	smtpClient             *smtp.Client
	smsClient              *twilio.RestClient
}

func NewNotification(notificationRepository ports.NotificationRepository, redisClient *redis.Client, smtpClient *smtp.Client, smsClient *twilio.RestClient) *notificationService {
	NotificationService = &notificationService{
		notificationRepository: notificationRepository,
		redisClient:            redisClient,
		smtpClient:             smtpClient,
		smsClient:              smsClient,
	}
	return NotificationService
}

func (service *notificationService) CreateNotification(createNotificationDto dto.CreateNotification) (interface{}, error) {
	logger.LogEvent("INFO", "Creating Notification")

	notification := mapper.MapCreateDto(createNotificationDto)

	notificationCreatedEvent := event.NotificationCreatedEvent{
		Event: eto.Event{
			EventReference:     uuid.New().String(),
			EventName:          "notificationcreatedevent",
			EventDate:          time.Now().Format(time.RFC3339),
			EventType:          "notificationcreatedevent",
			EventSource:        configuration.ServiceConfiguration.ServiceName,
			EventUserReference: createNotificationDto.UserReference,
			EventData:          notification,
		},
	}
	//saves to database
	response, err := service.notificationRepository.CreateNotification(notification)

	if err != nil {
		return nil, err
	}

	if notification.Channel == shared.Email {
		//send to gmail smtp
		emailSender := sender.NewEmailNotification(service.smtpClient)
		emailSender.SendEmailNotification(notification)
	} else if notification.Channel == shared.Sms {
		smsSender := sender.NewSmsNotification(service.smsClient)
		smsSender.SendSmsNotification(notification)
	}

	//publish to redis channel
	eventPublisher := publisher.NewPublisher(service.redisClient)
	ctx := context.Background()
	eventPublisher.PublishToSignUpCreatedEvent(ctx, notificationCreatedEvent)

	return response, err

}

func (service *notificationService) GetNotificationByDeviceReference(device_reference string, page string) (interface{}, error) {
	logger.LogEvent("INFO", "Getting the notification list for device "+device_reference)
	notificationlist, err := service.notificationRepository.GetNotificationByDeviceReference(device_reference, page)

	if err != nil {
		return nil, err
	}
	return notificationlist, nil
}

func (service *notificationService) GetNotificationByReference(reference string) (interface{}, error) {
	logger.LogEvent("INFO", "Getting notification with reference: "+reference)
	notification, err := service.notificationRepository.GetNotificationByReference(reference)

	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (service *notificationService) GetNotificationByUserReference(user_reference string) (interface{}, error) {
	logger.LogEvent("INFO", "Getting notification by user reference: "+user_reference)
	notification, err := service.notificationRepository.GetNotificationByUserReference(user_reference)

	if err != nil {
		return nil, err
	}
	return notification, nil
}
