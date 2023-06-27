package services

import (
	"context"

	"time"
	publisher "walls-notification-service/internal/adapter/events/publisher"
	messenger "walls-notification-service/internal/adapter/twilio-events"
	"walls-notification-service/internal/core/domain/dto"
	event "walls-notification-service/internal/core/domain/event/eto"
	"walls-notification-service/internal/core/domain/mapper"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
	eto "walls-notification-service/internal/core/helper/event-helper/eto"
	logger "walls-notification-service/internal/core/helper/log-helper"

	ports "walls-notification-service/internal/port"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sendgrid/sendgrid-go"
	"github.com/twilio/twilio-go"
)

var NotificationService = &notificationService{}

type notificationService struct {
	notificationRepository ports.NotificationRepository
	redisClient            *redis.Client
	smsClient              *twilio.RestClient
	emailClient            *sendgrid.Client
}

func NewNotification(notificationRepository ports.NotificationRepository, redisClient *redis.Client,
	smsClient *twilio.RestClient, emailClient *sendgrid.Client) *notificationService {
	NotificationService = &notificationService{
		notificationRepository: notificationRepository,
		redisClient:            redisClient,
		smsClient:              smsClient,
		emailClient:            emailClient,
	}
	return NotificationService
}

func (service *notificationService) CreateNotification(createNotificationDto dto.CreateNotification) (interface{}, error) {
	logger.LogEvent("INFO", "Creating Notification")

	mappedNotification := mapper.MapCreateDto(createNotificationDto)

	notificationCreatedEvent := event.NotificationCreatedEvent{
		Event: eto.Event{
			EventReference:     uuid.New().String(),
			EventName:          "notificationcreatedevent",
			EventDate:          time.Now().Format(time.RFC3339),
			EventType:          "notificationcreatedevent",
			EventSource:        configuration.ServiceConfiguration.ServiceName,
			EventUserReference: createNotificationDto.UserReference,
			EventData:          mappedNotification,
		},
	}

	//saves to database
	response, err := service.notificationRepository.CreateNotification(mappedNotification)

	if err != nil {
		return nil, err
	}

	// send sms using twilio or email using sendgrid
	messenger := messenger.NewMessenger(service.smsClient, service.emailClient)
	if err := messenger.SendNotificationMessage(mappedNotification, service.notificationRepository.UpdateNotifcation); err != nil {
		return nil, err
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

func (service *notificationService) GetNotificationByUserReference(user_reference string, page string) (interface{}, error) {
	logger.LogEvent("INFO", "Getting notification by user reference: "+user_reference)
	notification, err := service.notificationRepository.GetNotificationByUserReference(user_reference, page)

	if err != nil {
		return nil, err
	}
	return notification, nil
}
