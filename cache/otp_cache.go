package cache

import (
	"fmt"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/redis/go-redis/v9"
)

func GetOtpFromCache(key string) (string, error) {
	data, err := initializers.RedisClient.Get(ctx, "otp-"+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("OTP not found in cache")
		}
		go helpers.LogServerError("Error getting OTP from cache", err, "redis")
		return "", fmt.Errorf("error getting OTP from cache")
	}
	return data, nil
}

func SetOtpToCache(key string, data []byte) error {
	if err := initializers.RedisClient.Set(ctx, "otp-"+key, data, config.VERIFICATION_OTP_EXPIRATION_TIME).Err(); err != nil {
		go helpers.LogServerError("Error setting OTP to cache", err, "redis")
		return fmt.Errorf("error setting OTP to cache")
	}
	return nil
}
