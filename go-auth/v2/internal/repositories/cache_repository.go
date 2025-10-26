package repositories

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepositoryInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) CacheRepositoryInterface {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *cacheRepository) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
