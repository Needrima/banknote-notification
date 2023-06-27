package messenger

import (
	"walls-notification-service/internal/core/domain/entity"
	helper "walls-notification-service/internal/core/helper/twilio-helper"

	"github.com/sendgrid/sendgrid-go"
	"github.com/twilio/twilio-go"
)

type Messenger struct {
	smsClient   *twilio.RestClient
	emailClient *sendgrid.Client
}

func NewMessenger(smsClient *twilio.RestClient, emailClient *sendgrid.Client) *Messenger {
	return &Messenger{
		smsClient:   smsClient,
		emailClient: emailClient,
	}
}

// SendNotificationMessage sends a message based on notification.Channel and notification.Type
// and calls updateFunc to update notification.Status and notification.NotifiedOn fileds in the database if the message was sent successfully.
func (m *Messenger) SendNotificationMessage(notification entity.Notification, updateFunc func(reference string) error) error {
	client := helper.NewMessengerClient(m.smsClient, m.emailClient)
	return client.SendMessage(notification, updateFunc)
}
