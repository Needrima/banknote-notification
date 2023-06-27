package helper

import (
	"encoding/json"
	"fmt"
	"time"
	"walls-notification-service/internal/core/domain/entity"
	"walls-notification-service/internal/core/domain/shared"
	config "walls-notification-service/internal/core/helper/configuration-helper"
	logger "walls-notification-service/internal/core/helper/log-helper"
	timeHelper "walls-notification-service/internal/core/helper/parse-time-helper"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessengerClient struct {
	smsClient   *twilio.RestClient
	emailClient *sendgrid.Client
}

func NewMessengerClient(smsClient *twilio.RestClient, emailClient *sendgrid.Client) *MessengerClient {
	return &MessengerClient{
		smsClient:   smsClient,
		emailClient: emailClient,
	}
}

// SendMessage sends a message based on notification.Channel and notification.Type
// and calls updateFunc to update notification.Status and notification.NotifiedOn fileds in the database if the message was sent successfully.
func (r *MessengerClient) SendMessage(notification entity.Notification, updateFunc func(reference string) error) error {
	switch notification.Channel {
	case shared.Email:

		switch notification.Type {
		case shared.Instant:
			err := sendEmail(r.emailClient, notification)
			if err == nil {
				updateFunc(notification.Reference)
			}

			return err

		case shared.Scheduled: // launch a go routine to wait until scheduled time before sending mail
			timeToSendNotification := timeHelper.PeriodToScheduledTime(notification.NotifyOn)
			if timeToSendNotification < 10 { // time to send notification must be atleast 10 seconds later since mail is scheduled
				logger.LogEvent("INFO", "time to send scheduled mail not a future time")
				return fmt.Errorf("time to send scheduled mail not a future time")
			}

			ticker := time.NewTicker(time.Second * time.Duration(timeToSendNotification))
			go func() {
				<-ticker.C
				err := sendEmail(r.emailClient, notification)
				if err != nil {
					logger.LogEvent("ERROR", err.Error())
					return
				}
				updateFunc(notification.Reference)
			}()

		default:
			return fmt.Errorf("notification type %v is not allowed", notification.Type)
		}

	case shared.Sms:

		switch notification.Type {
		case shared.Instant:
			err := sendSms(r.smsClient, notification)
			if err == nil {
				updateFunc(notification.Reference)
			}

			return err

		case shared.Scheduled: // launch a go routine to wait until scheduled time before sending mail
			timeToSendNotification := timeHelper.PeriodToScheduledTime(notification.NotifyOn)
			if timeToSendNotification < 10 { // time to send notification must be atleast 10 seconds later since mail is scheduled
				logger.LogEvent("INFO", "time to send scheduled sms not a future time")
				return fmt.Errorf("time to send scheduled sms not a future time")
			}

			ticker := time.NewTicker(time.Second * time.Duration(timeToSendNotification))
			go func() {
				<-ticker.C
				err := sendSms(r.smsClient, notification)
				if err != nil {
					logger.LogEvent("ERROR", err.Error())
					return
				}
				updateFunc(notification.Reference)
			}()

		default:
			return fmt.Errorf("notification type %v is not allowed", notification.Type)
		}
	}

	return fmt.Errorf("notification channel %v is not allowed", notification.Channel)
}

// sendEmail send email using twilio sendgrid api
func sendEmail(emailClient *sendgrid.Client, notification entity.Notification) error {

	// Create the email message
	from := mail.NewEmail("Wallspay", "wallspay@sendgrid.com")
	subject := notification.Subject
	to := mail.NewEmail("", notification.Contact)
	content := mail.NewContent("text/plain", notification.MessageBody)
	message := mail.NewV3MailInit(from, subject, to, content)

	// Send the email
	response, err := emailClient.Send(message)                    
	if err != nil {
		logger.LogEvent("ERROR", fmt.Sprintf("Error sending email: %s, reference: %s", err.Error(), notification.Reference))
		return err
	}

	if response.StatusCode < 200 && response.StatusCode >= 300 {
		logger.LogEvent("ERORR", fmt.Sprintf("Error sending email, returned status code: %d, reference: %s", response.StatusCode, notification.Reference))
		return fmt.Errorf("error sending email with status code: %d", response.StatusCode)
	}

	resp, _ := json.Marshal(*response)
	logger.LogEvent("INFO", fmt.Sprintf("Sending email successful reference: %s, message: %s", notification.Reference, resp))

	return nil
}

// sendEmail send sms using twilio sms api
func sendSms(client *twilio.RestClient, notification entity.Notification) error {
	params := &twilioApi.CreateMessageParams{
		To:   &notification.Contact,
		From: &config.ServiceConfiguration.TwilioAuthPhoneNumber,
		Body: &notification.MessageBody,
	}

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		logger.LogEvent("ERROR", fmt.Sprintf("Error sending SMS message: %s", err.Error()))
		return err
	}

	response, _ := json.Marshal(*resp)
	logger.LogEvent("INFO", fmt.Sprintf("Sending sms successful reference: %s, message: %s", notification.Reference, response))

	return nil
}
