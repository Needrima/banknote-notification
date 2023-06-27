package extensions

import (
	"strings"
	messaging "walls-notification-service/internal/adapter/twiliomsg"
	logger "walls-notification-service/internal/core/helper/log-helper"

	"github.com/sendgrid/sendgrid-go"
	"github.com/twilio/twilio-go"
)

func StartTwilioConnection(connType string) *twilio.RestClient {

	switch connType {
	case strings.ToLower(connType):
		logger.LogEvent("INFO", "Connectiong to twilio!")
		smsClient := messaging.ConnectToTwilio()
		return smsClient
	}
	return nil

}

func StartSendGridConnection(connType string) *sendgrid.Client {

	switch connType {
	case strings.ToLower(connType):
		logger.LogEvent("INFO", "Connectiong to twilio!")
		emailClient := messaging.ConnectToSendgrid()
		return emailClient
	}
	return nil

}
