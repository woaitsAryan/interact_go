package cache

import (
	"context"
	"fmt"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/redis/go-redis/v9"
)

var ctx = context.TODO()

func GetFromCache(key string) (string, error) {
	data, err := initializers.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("item not found in cache")
		}
		helpers.LogServerError("Error Getting from cache", err, "")
		return "", fmt.Errorf("error getting from cache")
	}
	return data, nil
}

func SetToCache(key string, data []byte) error {
	if err := initializers.RedisClient.Set(ctx, key, data, initializers.CacheExpirationTime).Err(); err != nil {
		helpers.LogServerError("Error Setting to cache", err, "")
		return fmt.Errorf("error setting to cache")
	}
	return nil
}
