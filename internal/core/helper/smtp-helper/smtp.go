package helper

import (
	"net/smtp"
	"walls-notification-service/internal/core/domain/entity"
)

type SmtpClient struct {
	client *smtp.Client
}

func NewSmtpClient(client *smtp.Client) *SmtpClient {
	return &SmtpClient{
		client: client,
	}
}

func (s *SmtpClient) SendEmail(notification entity.Notification) error {
	// Set the sender and recipient
	if err := s.client.Mail(notification.From); err != nil {
		return err
	}
	if err := s.client.Rcpt(notification.Contact); err != nil {
		return err
	}

	// Send the email
	w, err := s.client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(notification.MessageBody))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
