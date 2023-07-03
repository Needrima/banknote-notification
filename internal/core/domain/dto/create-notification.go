package dto

import "walls-notification-service/internal/core/domain/shared"

type CreateNotification struct {
	UserReference 	string 						`json:"user_reference" bson:"user_reference" validate:"required,guid,min=32,max=38"`
	DeviceReference string 						`json:"device_reference" bson:"device_reference" validate:"required,omitempty,guid,min=32,max=38"`
	Contact       	string       				`json:"contact" bson:"contact" validate:"required,valid_contact"`
	Channel       	shared.Channel  			`json:"channel" bson:"channel" validate:"required,valid_channel"`
	Subject      	string						`json:"subject" bson:"subject" validate:"required"`
	MessageBody     string            			`json:"message_body" bson:"message_body" validate:"required"`
}
