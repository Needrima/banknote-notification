package repository

import (
	configuration "walls-notification-service/internal/core/helper/configuration-helper"

	"github.com/twilio/twilio-go"
)

func ConnectToTwilioSms() *twilio.RestClient {
	config := configuration.ServiceConfiguration
	accountSid := config.TwilioAccountSID
	authToken := config.TwilioAuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return client
}
