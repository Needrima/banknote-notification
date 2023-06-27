package dto

import "walls-notification-service/internal/core/domain/shared"

type CreateNotification struct {
	UserReference   string                  `json:"user_reference" bson:"user_reference" validate:"required,min=26,max=38"`
	DeviceReference string                  `json:"device_reference" bson:"device_reference"`
	Contact         string                  `json:"contact" bson:"contact" validate:"required,valid_contact"`
	Channel         shared.Channel          `json:"channel" bson:"channel" validate:"required,valid_channel"`
	Type            shared.NotificationType `json:"notification_type" bson:"notification_type" validate:"required"`
	Subject         string                  `json:"subject" bson:"subject" validate:"required"`
	MessageBody     string                  `json:"message_body" bson:"message_body" validate:"required"`
	NotifiedBy      string                  `json:"notified_by" bson:"notified_by" validate:"required,min=26,max=38"`
	NotifyOn        string                  `json:"notify_on" bson:"notify_on" validate:"required"`
}
