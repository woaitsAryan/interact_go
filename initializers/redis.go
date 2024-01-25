package initializers

import (
	"context"
	"fmt"
	"time"

	"github.com/Pratham-Mishra04/interact/cache/subscribers"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

var ctx = context.TODO()

var CacheExpirationTime = time.Second * 10
var CacheExpirationTimeLong = time.Hour * 24

func ConnectToCache() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     CONFIG.REDIS_HOST + ":" + CONFIG.REDIS_PORT,
		Password: CONFIG.REDIS_PASSWORD,
		DB:       0,
		PoolSize: 1000,
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("Redis connection Error:\n %v", err)
	} else {
		fmt.Println("Connected to redis!")

		RedisExpirationSub := RedisClient.Subscribe(ctx, "__keyevent@0__:expired")
		defer RedisExpirationSub.Close()

		// Wait for confirmation that subscription is created before publishing anything
		_, err := RedisExpirationSub.Receive(ctx)
		if err != nil {
			fmt.Println("Error Subscribing to Redis Expiration Event: ", err)
		} else {
			fmt.Println("Subscribed to Redis Expiration Event")
			go subscribers.ImpressionsDumpSub(RedisExpirationSub, DB)
		}
	}
}
