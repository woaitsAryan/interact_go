package initializers

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

var ctx = context.TODO()

var CacheExpirationTime = time.Minute * 10

func ConnectToCache() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection Error: %v", err.Error())
	}
}
