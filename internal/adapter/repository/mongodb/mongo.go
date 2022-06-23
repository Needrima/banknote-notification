package repository

import (
	"bnt/bnt-notification-service/internal/core/helper"
	ports "bnt/bnt-notification-service/internal/port"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type MongoRepositories struct {
	Notification ports.NotificationRepository
}

func ConnectToMongo(dbType string, dbUsername string, dbPassword string, dbHost string,
	dbPort string, authdb string, dbname string) (MongoRepositories, error) {
	helper.LogEvent("INFO", "Establishing mongoDB connection with given credentials...")
	var mongoCredentials, authSource string
	if dbUsername != "" && dbPassword != "" {
		mongoCredentials = fmt.Sprint(dbUsername, ":", dbPassword, "@")
		authSource = fmt.Sprint("authSource=", authdb, "&")
	}
	mongoUrl := fmt.Sprint(dbType, "://", mongoCredentials, dbHost, ":", dbPort, "/?", authSource,
		"directConnection=true&serverSelectionTimeoutMS=2000")
	clientOptions := options.Client().ApplyURI(mongoUrl) // Connect to
	helper.LogEvent("INFO", "Connecting to MongoDB...")
	db, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		//log.Println(err)
		//log.Fatal(err)
		helper.LogEvent("ERROR", helper.ErrorMessage(helper.MongoDBError, err.Error()))
		return MongoRepositories{}, err
	}

	// Check the connection
	helper.LogEvent("INFO", "Confirming MongoDB Connection...")
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		//log.Println(err)
		//log.Fatal(err)
		helper.LogEvent("ERROR", helper.ErrorMessage(helper.MongoDBError, err.Error()))
		return MongoRepositories{}, err
	}

	//helper.LogEvent("Info", "Connected to MongoDB!")
	helper.LogEvent("INFO", "Establishing Database collections and indexes...")
	conn := db.Database(dbname)

	notificationCollection := conn.Collection("notifications")

	repo := MongoRepositories{
		Notification: NewNotification(notificationCollection),
	}
	return repo, nil
}

func GetPage(page string) (*options.FindOptions, error) {
	if page == "all" {
		return nil, nil
	}
	var limit, e = strconv.ParseInt(helper.Config.PageLimit, 10, 64)
	var pageSize, ee = strconv.ParseInt(page, 10, 64)
	if e != nil || ee != nil {
		return nil, helper.ErrorMessage(helper.NoRecordError, "Error in page-size or limit-size.")
	}
	findOptions := options.Find().SetLimit(limit).SetSkip(limit * (pageSize - 1))
	return findOptions, nil
}
