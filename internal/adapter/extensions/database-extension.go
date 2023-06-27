package extensions

import (
	mongoRepository "walls-notification-service/internal/adapter/repository/mongodb"
	logger "walls-notification-service/internal/core/helper/log-helper"

	"fmt"
	"log"
	"strings"
)

func StartDatabase(dbType string) interface{} {

	switch dbType {
	case strings.ToLower(dbType):
		logger.LogEvent("INFO", "Initializing Mongo!")
		mongoRepo, err := mongoRepository.ConnectToMongo()
		if err != nil {
			fmt.Println(err)
			logger.LogEvent("ERROR", "MongoDB database Connection Error: "+err.Error())
			log.Fatal()
		}

		return mongoRepo
	}
	return nil

}
