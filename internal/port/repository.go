package ports

import (
	"walls-notification-service/internal/core/domain/entity"
)

type NotificationRepository interface {
	CreateNotification(createNotificationDto entity.Notification) (interface{}, error)
	GetNotificationByReference(reference string) (interface{}, error)
	GetNotificationByDeviceReference(device_reference string, page string) (interface{}, error)
	GetNotificationByUserReference(user_reference string) (interface{}, error)
}
