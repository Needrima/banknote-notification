package event

import (
	"walls-notification-service/internal/core/helper/event-helper/eto"
)

type NotificationCreatedEvent struct {
	eto.Event
}
