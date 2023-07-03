package extensions

import (
	"net/smtp"
	"strings"
	google "walls-notification-service/internal/adapter/repository/google"
	logger "walls-notification-service/internal/core/helper/log-helper"


)

func StartSmtp(smtpType string) *smtp.Client {

	switch smtpType {
	case strings.ToLower(smtpType):
		logger.LogEvent("INFO", "Initializing Smtp!")
		smtpConnection,err := google.ConnectToGoogleSmtp()
		if (err != nil) {
			return &smtp.Client{}
		}
		return smtpConnection
	}
	return nil

}
