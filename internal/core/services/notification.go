package services

import (
	"bnt/bnt-notification-service/internal/core/domain/entity"
	"bnt/bnt-notification-service/internal/core/helper"
	ports "bnt/bnt-notification-service/internal/port"
	"github.com/google/uuid"
)

type notificationService struct {
	notificationRepository ports.NotificationRepository
}

func NewNotification(notificationRepository ports.NotificationRepository) *notificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}

func (service *notificationService) GetNotificationList(page string) (interface{}, error) {
	helper.LogEvent("INFO", "Getting all categories...")
	notification, err := service.notificationRepository.GetNotificationList(page)

	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (service *notificationService) CreateNotification(notification entity.Notification) (interface{}, error) {
	notification.Reference = uuid.New().String()
	helper.LogEvent("INFO", "Creating notification with reference: "+notification.Reference)
	// if err := helper.Validate(notification); err != nil {
	// 	return nil, err
	// }
	return service.notificationRepository.CreateNotification(notification)
}

func (service *notificationService) GetNotificationStatus(reference string) (interface{}, error) {
	helper.LogEvent("INFO", "Enabling notification with reference: "+reference)
	_, err := service.GetNotificationByRef(reference)
	if err != nil {
		return nil, err
	}
	return service.notificationRepository.GetNotificationStatus(reference)
}

func (service *notificationService) GetNotificationByRef(reference string) (interface{}, error) {
	helper.LogEvent("INFO", "Getting notification with reference: "+reference)
	notification, err := service.notificationRepository.GetNotificationByRef(reference)
	if err != nil {
		return nil, err
	}
	return notification, nil
}
