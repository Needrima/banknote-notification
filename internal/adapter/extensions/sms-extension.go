package extensions

import (
	"strings"
	twilioRepo "walls-notification-service/internal/adapter/repository/twilio"
	logger "walls-notification-service/internal/core/helper/log-helper"

	"github.com/twilio/twilio-go"
)

func StartSms(smsType string) *twilio.RestClient {

	switch smsType {
	case strings.ToLower(smsType):
		logger.LogEvent("INFO", "Initializing Sms!")
		smsConnection := twilioRepo.ConnectToTwilioSms()
		return smsConnection
	}
	return nil

}
