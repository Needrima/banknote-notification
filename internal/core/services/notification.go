package services

import (
	"fmt"
	"github.com/google/uuid"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/helper"
	ports "notification-service/internal/port"
	"time"
)

type notificationService struct {
	notificationRepository ports.NotificationRepository
}

func NewNotification(notificationRepository ports.NotificationRepository) *notificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}

func (service *notificationService) GetNotificationList(page string) (interface{}, error) {
	helper.LogEvent("INFO", "Getting all categories...")
	notification, err := service.notificationRepository.GetNotificationList(page)

	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (service *notificationService) CreateNotification(notification entity.Notification) (interface{}, error) {
	notification.Reference = uuid.New().String()
	switch notification.Type {
	case "INSTANT":
		// send mail immediately
		if err := helper.SendMail(notification.To, notification.From, notification.Message); err != nil {
			helper.LogEvent("INFO", "sending mail for instant notification failed")
			return nil, helper.ErrorMessage(helper.ServerError, "something went wrong")
		}

		sendTime := time.Now()
		notification.SentAt = helper.ParseTimeToString(sendTime)
		notification.Status = "SENT"

		return service.notificationRepository.CreateNotification(notification)

	case "SCHEDULED":
		// get the scheduled time for sending the notification and parse to time variable
		scheduledTime, err := helper.ParseTimeStringToTime(notification.SendAt)
		if err != nil {
			helper.LogEvent("INFO", "invalid time format")
			return nil, helper.ErrorMessage(helper.ValidationError, "send_at time not a valid time format")
		}

		// get time to scheduled time in seconds
		sendTime := helper.PeriodToScheduledTime(scheduledTime)

		// check if time to scheduled time is a future time (i.e at least 10 seconds later)
		if sendTime < 10 {
			helper.LogEvent("INFO", "time to send scheduled mail not a future time")
			return nil, helper.ErrorMessage(helper.InvalidScheduleDate, "You cannot schedule a task in the past. You must provide a future date")
		}

		// set status to pending
		notification.Status = "PENDING"

		// launch a goroutine to check if scheduled time has elapsed and
		// update the database once the schedule time reaches
		ticker := time.NewTicker(time.Second * time.Duration(sendTime))
		go func() {
			select {
			case <-ticker.C:
				// send notification mail
				if err := helper.SendMail(notification.To, notification.From, notification.Message); err != nil {
					helper.LogEvent("INFO", fmt.Sprintf("sending mail for scheduled notification with reference %v: %v", notification.Reference, err.Error()))
					return
				}

				// update notification send_at time and status from "PENDING" to "SENT"
				notification.SentAt = helper.ParseTimeToString(time.Now())
				notification.Status = "SENT"

				// update notification in database
				ref, _ := service.notificationRepository.UpdateNotification(notification)
				helper.LogEvent("INFO", fmt.Sprintf("updating scheduled notification with reference %v", ref))
			}
		}()

		// insert notification to database
		return service.notificationRepository.CreateNotification(notification)
	default:
		return "", nil
	}
}

func (service *notificationService) GetNotificationStatus(reference string) (interface{}, error) {
	helper.LogEvent("INFO", "Enabling notification with reference: "+reference)
	_, err := service.GetNotificationByRef(reference)
	if err != nil {
		return nil, err
	}
	return service.notificationRepository.GetNotificationStatus(reference)
}

func (service *notificationService) GetNotificationByRef(reference string) (interface{}, error) {
	helper.LogEvent("INFO", "Getting notification with reference: "+reference)
	notification, err := service.notificationRepository.GetNotificationByRef(reference)
	if err != nil {
		return nil, err
	}
	return notification, nil
}
