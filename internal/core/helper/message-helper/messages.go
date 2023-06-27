package helper

import (
	"fmt"
	"time"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
)

var (
	NoResourceFound = "this resource does not exist"
	NoRecordFound   = "sorry. no record found"
	ServerStarted   = fmt.Sprintf("Started " + configuration.ServiceConfiguration.ServiceName + " on " + configuration.ServiceConfiguration.ServiceAddress + ":" +
		configuration.ServiceConfiguration.ServicePort + " in " + time.Since(time.Now()).String())
	StartingServer = "Attempting to start " + configuration.ServiceConfiguration.ServiceName
	StartingRedis  = "Attempting to start Redis on " + configuration.ServiceConfiguration.EBConnectionString
)

var (
	ValidationError = "VALIDATION_ERROR"
	RedisSetupError = "REDIS_SETUP_ERROR"
	NoRecordError   = "NO_RECORD_FOUND_ERROR"
	NoResourceError = "INVALID_RESOURCE_ERROR"
	CreateError     = "CREATE_ERROR"
	UpdateError     = "UPDATE_ERROR"
	LogError        = "LOG_ERROR"
	MongoDBError    = "MONGO_DB_ERROR"
	RedisError      = "REDIS_ERROR"
)
