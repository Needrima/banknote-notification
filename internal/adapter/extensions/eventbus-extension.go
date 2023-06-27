package extensions

import (
	"strings"
	redis "walls-notification-service/internal/adapter/repository/redis"
	logger "walls-notification-service/internal/core/helper/log-helper"

	redisClient "github.com/redis/go-redis/v9"
)

func StartEventBus(ebType string) *redisClient.Client {

	switch ebType {
	case strings.ToLower(ebType):
		logger.LogEvent("INFO", "Initializing Redis!")
		redisRepo := redis.ConnectToRedis()
		return redisRepo
	}
	return nil

}
