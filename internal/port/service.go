package ports

import (
	"walls-notification-service/internal/core/domain/dto"
)

type NotificationService interface {
	CreateNotification(createNotificationDto dto.CreateNotification) (interface{}, error)
	GetNotificationByReference(reference string) (interface{}, error)
	GetNotificationByDeviceReference(device_reference string, page string) (interface{}, error)
	GetNotificationByUserReference(user_reference string, page string) (interface{}, error)
}
