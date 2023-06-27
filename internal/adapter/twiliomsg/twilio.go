package messaging

import (
	helper "walls-notification-service/internal/core/helper/configuration-helper"

	"github.com/sendgrid/sendgrid-go"
	"github.com/twilio/twilio-go"
)

func ConnectToTwilio() *twilio.RestClient {
	accountSid := helper.ServiceConfiguration.TwilioAccountSID
	authToken := helper.ServiceConfiguration.TwilioAuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return client
}

func ConnectToSendgrid() *sendgrid.Client {
	apiKey := helper.ServiceConfiguration.SendgridAPIKey
	client := sendgrid.NewSendClient(apiKey)

	return client
}
