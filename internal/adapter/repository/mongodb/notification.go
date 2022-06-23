package repository

import (
	"bnt/bnt-notification-service/internal/core/domain/entity"
	"bnt/bnt-notification-service/internal/core/helper"
	ports "bnt/bnt-notification-service/internal/port"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"reflect"
)

type NotificationInfra struct {
	Collection *mongo.Collection
}

func NewNotification(Collection *mongo.Collection) *NotificationInfra {
	return &NotificationInfra{Collection}
}

//UserRepo implements the repository.UserRepository interface
var _ ports.NotificationRepository = &NotificationInfra{}

func (r *NotificationInfra) CreateNotification(notification entity.Notification) (interface{}, error) {
	helper.LogEvent("INFO", "Persisting notification configurations with reference: "+notification.Reference)
	switch notification.Type {
	case "INSTANT":
		// send mail immediately
		if err := helper.SendMail(notification.From, notification.To, notification.Message); err != nil {
			helper.LogEvent("INFO", "sending mail for instant notification failed")
			return nil, helper.ErrorMessage(helper.ServerError, "something went wrong")
		}

		sendTime := time.Now()
		notification.SentAt = helper.ParseTimeToString(sendTime)
		notification.Status = "SENT"

		// insert notification into the database
		_, err := r.Collection.InsertOne(context.TODO(), notification)
		if err != nil {
			helper.LogEvent("INFO", "error inserting document into mogodb "+err.Error())
			return nil, helper.ErrorMessage(helper.MongoDBError, "something went wrong")
		}

		// send the notification reference as a response
		return notification.Reference, nil

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

		// insert notification to database
		_, err = r.Collection.InsertOne(context.TODO(), notification)
		if err != nil {
			helper.LogEvent("INFO", "error inserting document into mogodb "+err.Error())
			return nil, helper.ErrorMessage(helper.MongoDBError, "something went wrong")
		}

		// set a ticker and launch a goroutine to check if scheduled time has elapsed and launch a
		// update the database once the schedule time reaches
		ticker := time.NewTicker(time.Second * time.Duration(sendTime))
		go func(collection *mongo.Collection, notification entity.Notification) {
			select {
			case <-ticker.C:
				// send notification mail
				if err := helper.SendMail(notification.From, notification.To, notification.Message); err != nil {
					helper.LogEvent("INFO", fmt.Sprintf("sending mail for scheduled notification with reference %v: %v", notification.Reference, err.Error()))
					return
				}

				// update notification send_at time and status from "PENDING" to "SENT"
				notification.SentAt = helper.ParseTimeToString(time.Now())
				notification.Status = "SENT"

				// update notification in database
				_, err := collection.UpdateOne(context.TODO(), bson.M{"reference": notification.Reference}, bson.M{"$set": notification})
				if err != nil {
					helper.LogEvent("INFO", fmt.Sprintf("updating scheduled notification with reference %v: %v", notification.Reference, err.Error()))
					return
				}
			}
		}(r.Collection, notification)
	}

	return notification.Reference, nil
}

func (r *NotificationInfra) GetNotificationStatus(reference string) (interface{}, error) {
	helper.LogEvent("INFO", "Retrieving notification configurations with reference: "+reference)
	notification := entity.Notification{}
	filter := bson.M{"reference": reference}
	err := r.Collection.FindOne(context.TODO(), filter).Decode(&notification)
	if err != nil || notification == (entity.Notification{}) {
		return nil, helper.ErrorMessage(helper.NoRecordError, helper.NoRecordFound)
	}
	helper.LogEvent("INFO", "Retrieving notification configurations with reference: "+reference+" completed successfully. ")
	return notification.Status, nil
}

func (r *NotificationInfra) GetNotificationByRef(reference string) (interface{}, error) {
	helper.LogEvent("INFO", "Retrieving notification configurations with reference: "+reference)
	notification := entity.Notification{}
	filter := bson.M{"reference": reference}
	err := r.Collection.FindOne(context.TODO(), filter).Decode(&notification)
	if err != nil || notification == (entity.Notification{}) {
		return nil, helper.ErrorMessage(helper.NoRecordError, helper.NoRecordFound)
	}
	helper.LogEvent("INFO", "Retrieving notification configurations with reference: "+reference+" completed successfully. ")
	return notification, nil
}

func (r *NotificationInfra) GetNotificationList(page string) (interface{}, error) {
	helper.LogEvent("INFO", "Retrieving all notification configuration entries...")
	var notifications []entity.Notification
	var notification entity.Notification
	findOptions, err := GetPage(page)
	if err != nil {
		return nil, helper.ErrorMessage(helper.NoRecordError, "Error in page-size or limit-size.")
	}
	cursor, err := r.Collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		return nil, helper.ErrorMessage(helper.NoRecordError, helper.NoRecordFound)
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&notification)
		if err != nil {

			return nil, helper.ErrorMessage(helper.NoRecordError, err.Error())
		}
		notifications = append(notifications, notification)
	}
	if reflect.ValueOf(notifications).IsNil() {
		helper.LogEvent("INFO", "There are no results in this collection...")
		return []entity.Notification{}, nil
	}
	helper.LogEvent("INFO", "Retrieving all notification configuration entries completed successfully")
	return notifications, nil
}
