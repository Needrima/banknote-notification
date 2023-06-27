package publisher

import (
	"context"
	helper "walls-notification-service/internal/core/helper/event-helper"

	"github.com/redis/go-redis/v9"
)

type EventPublisher struct {
	redisClient *redis.Client
}

func NewPublisher(redisClient *redis.Client) *EventPublisher {
	return &EventPublisher{
		redisClient: redisClient,
	}
}

func (p *EventPublisher) PublishToSignUpCreatedEvent(ctx context.Context, event interface{}) error {

	redisHelper := helper.NewRedisClient(p.redisClient)
	return redisHelper.PublishEvent(ctx, event)

}
