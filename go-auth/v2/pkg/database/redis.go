package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisInstance struct {
	Client *redis.Client
}

var RedisDB RedisInstance

func ConnectRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	log.Println("Redis connected successfully")
	RedisDB = RedisInstance{Client: rdb}
}

func GetRedis() *redis.Client {
	return RedisDB.Client
}
