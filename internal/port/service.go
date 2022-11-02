package ports

import "notification-service/internal/core/domain/entity"

type NotificationService interface {
	CreateNotification(notification entity.Notification) (interface{}, error)
	GetNotificationStatus(ref string) (interface{}, error)
	GetNotificationByRef(ref string) (interface{}, error)
	GetNotificationList(page string) (interface{}, error)
}
