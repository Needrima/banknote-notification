package main

import (
	mongoRepository "notification-service/internal/adapter/repository/mongodb"
	"notification-service/internal/adapter/routes"

	"fmt"
	"log"
	"notification-service/internal/core/helper"
)

func main() {
	//Initialize request Log
	helper.InitializeLog()
	//Start DB Connection
	mongoRepo := startMongo()
	helper.LogEvent("INFO", "MongoDB Initialized!")

	helper.LogEvent("INFO", "Redis Initialized!")
	//Set up routes
	router := routes.SetupRouter(mongoRepo.Notification)
	//Print custom message for server start

	fmt.Println(helper.ServerStarted)
	//Log server start event
	helper.LogEvent("INFO", helper.ServerStarted)
	//start server
	_ = router.Run(":" + helper.Config.ServicePort)
	//api.SetConfiguration
}

func startMongo() mongoRepository.MongoRepositories {
	helper.LogEvent("INFO", "Initializing Mongo!")
	mongoRepo, err := mongoRepository.ConnectToMongo(helper.Config.DbType, helper.Config.MongoDbUserName,
		helper.Config.MongoDbPassword, helper.Config.MongoDbHost,
		helper.Config.MongoDbPort, helper.Config.MongoDbAuthDb,
		helper.Config.MongoDbName)
	if err != nil {
		fmt.Println(err)
		helper.LogEvent("ERROR", "MongoDB database Connection Error: "+err.Error())
		log.Fatal()
	}
	return mongoRepo
}
