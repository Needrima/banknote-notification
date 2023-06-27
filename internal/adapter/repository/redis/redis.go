package repository

import (
	configuration "walls-notification-service/internal/core/helper/configuration-helper"

	logger "walls-notification-service/internal/core/helper/log-helper"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis() *redis.Client {
	logger.LogEvent("INFO", "Establishing redis connection with given credentials...")
	var redisClient = redis.NewClient(&redis.Options{
		Addr: configuration.ServiceConfiguration.EBConnectionString,
	})

	return redisClient
}
