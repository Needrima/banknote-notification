package sender

import (
	
	"net/smtp"
	"walls-notification-service/internal/core/domain/entity"
	helper "walls-notification-service/internal/core/helper/smtp-helper"

	
)

type EmailNotification struct {
	smtpClient *smtp.Client
}

func NewEmailNotification(smtpClient *smtp.Client) *EmailNotification {
	return &EmailNotification{
		smtpClient: smtpClient,
	}
}

func (e *EmailNotification) SendEmailNotification(notification entity.Notification) error {

	smtpHelper := helper.NewSmtpClient(e.smtpClient)
	return smtpHelper.SendEmail(notification)

}
