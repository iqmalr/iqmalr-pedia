package repositories

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type EventRepositoryInterface interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

type eventRepository struct {
	client *redis.Client
}

func NewEventRepository(client *redis.Client) EventRepositoryInterface {
	return &eventRepository{client: client}
}

func (r *eventRepository) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.client.Publish(ctx, channel, message).Err()
}
