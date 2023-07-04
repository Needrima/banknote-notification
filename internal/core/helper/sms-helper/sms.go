package helper

import (
	"encoding/json"
	"walls-notification-service/internal/core/domain/entity"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
	logger "walls-notification-service/internal/core/helper/log-helper"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SmsClient struct {
	client *twilio.RestClient
}

func NewSmsClient(client *twilio.RestClient) *SmsClient {
	return &SmsClient{
		client: client,
	}
}

func (s *SmsClient) SendSms(notification entity.Notification) error {
	// Set the sender and recipient
	config := configuration.ServiceConfiguration
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(notification.Contact)
	params.SetFrom(config.TwilioAuthPhoneNumber)
	params.SetBody(notification.MessageBody)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	response, _ := json.Marshal(*resp)
	logger.LogEvent("INFO", string(response))

	return nil
}
