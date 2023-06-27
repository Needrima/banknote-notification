package mapper

import (
	"time"
	"walls-notification-service/internal/core/domain/dto"
	"walls-notification-service/internal/core/domain/entity"

	"github.com/google/uuid"
)

func MapCreateDto(notificationDto dto.CreateNotification) entity.Notification {

	notificationMap := entity.Notification{
		Reference:     uuid.New().String(),
		UserReference: notificationDto.UserReference,
		Contact:       notificationDto.Contact,
		Channel:       notificationDto.Channel,
		Type:          notificationDto.Type,
		Subject:       notificationDto.Subject,
		MessageBody:   notificationDto.MessageBody,
		NotifiedBy:    notificationDto.NotifiedBy,
		NotifyOn:      time.Now().Format(time.RFC3339),
		NotifiedOn:    time.Now().Format(time.RFC3339),
	}
	return notificationMap
}
