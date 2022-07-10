package ports

import "bnt/bnt-notification-service/internal/core/domain/entity"

type NotificationRepository interface {
	CreateNotification(notification entity.Notification) (interface{}, error)
	GetNotificationStatus(ref string) (interface{}, error)
	GetNotificationByRef(CountryCode string) (interface{}, error)
	GetNotificationList(page string) (interface{}, error)
	UpdateNotification(notification entity.Notification) (interface{}, error)
}
