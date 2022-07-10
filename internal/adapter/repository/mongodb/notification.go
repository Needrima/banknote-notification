package repository

import (
	"bnt/bnt-notification-service/internal/core/domain/entity"
	"bnt/bnt-notification-service/internal/core/helper"
	ports "bnt/bnt-notification-service/internal/port"
	"context"
	// "fmt"
	// "time"

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
	// insert notification to database
	_, err := r.Collection.InsertOne(context.TODO(), notification)
	if err != nil {
		helper.LogEvent("INFO", "error inserting document into mogodb "+err.Error())
		return nil, helper.ErrorMessage(helper.MongoDBError, "something went wrong")
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

func (r *NotificationInfra) UpdateNotification(notification entity.Notification) (interface{}, error) {
	_, err := r.Collection.UpdateOne(context.TODO(), bson.M{"reference": notification.Reference}, bson.M{"$set": notification})
	if err != nil {
		helper.LogEvent("INFO", "error updating document in mogodb "+err.Error())
		return nil, helper.ErrorMessage(helper.MongoDBError, "something went wrong")
	}

	return notification.Reference, nil
}
