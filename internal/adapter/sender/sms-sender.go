package sender

import (
	"walls-notification-service/internal/core/domain/entity"
	smshelper "walls-notification-service/internal/core/helper/sms-helper"

	"github.com/twilio/twilio-go"
)

type SmsNotification struct {
	smsClient *twilio.RestClient
}

func NewSmsNotification(smsClient *twilio.RestClient) *SmsNotification {
	return &SmsNotification{
		smsClient: smsClient,
	}
}

func (e *SmsNotification) SendSmsNotification(notification entity.Notification) error {

	smsHelper := smshelper.NewSmsClient(e.smsClient)
	return smsHelper.SendSms(notification)

}
