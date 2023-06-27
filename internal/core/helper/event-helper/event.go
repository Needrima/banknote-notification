package helper

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strings"
	logger "walls-notification-service/internal/core/helper/log-helper"
	ports "walls-notification-service/internal/port"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) SubscribeToEvent(ctx context.Context, channel interface{}, eventHandler func(interface{}, ports.NotificationRepository)) error {
	// Get the channel name from the event object's type

	pubSub := r.client.Subscribe(ctx, channel.(string))
	defer pubSub.Close()

	ch := pubSub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			var event interface{}
			var port ports.NotificationRepository
			err := json.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Printf("Error decoding event: %v\n", err)
				logger.LogEvent("ERROR", "Error decoding event: "+err.Error())
				continue
			}

			eventHandler(event, port) // Pass the appropriate NotificationRepository instance here
		}
	}

}

func (r *RedisClient) PublishEvent(ctx context.Context, event interface{}) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Get the channel name from the event object's type
	eventChannel := strings.ToUpper(reflect.TypeOf(event).Name())

	err = r.client.Publish(ctx, eventChannel, string(eventBytes)).Err()
	if err != nil {
		return err
	}

	return nil
}
