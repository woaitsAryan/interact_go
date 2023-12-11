package cache

import (
	"fmt"
	"strconv"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/redis/go-redis/v9"
)

func SetImpression(key string, data int) error {
	if err := initializers.RedisClient.Set(ctx, "impression-" + key, data, initializers.CacheExpirationTimeLong).Err(); err != nil {
		helpers.LogServerError("Error Setting to impressions cache", err, "")
		return fmt.Errorf("error setting to impression cache")
	}
	return nil
}

func IncrementImpression(key string) error {
	data, err := initializers.RedisClient.Get(ctx, "impression-"+key).Result()
	if err != nil {
		if err == redis.Nil {
			if err := SetImpression(key, 1); err != nil {
				return err
			}
			return nil
		}
		helpers.LogServerError("Error Getting from impression cache", err, "")
		return fmt.Errorf("error getting from impression cache")
	}
	impressionCount, err := strconv.Atoi(data)
	if err != nil {
		helpers.LogServerError("Error converting impression count to int", err, "")
		return fmt.Errorf("error converting impression count to int")
	}
	if err := SetImpression(key, impressionCount+1); err != nil {
		return err
	}
	return nil
}

func GetImpression(key string) (int, error) {
	data, err := initializers.RedisClient.Get(ctx, "impression-"+key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		helpers.LogServerError("Error Getting from cache", err, "")
		return -1, fmt.Errorf("error getting from cache")
	}
	dataToInt, err := strconv.Atoi(data)
	if err != nil {
		helpers.LogServerError("Error converting impression count to int", err, "")
		return -1, fmt.Errorf("error converting impression count to int")
	}
	return dataToInt, nil
}

func ResetImpression(key string) error {
	if err := SetImpression(key, 0); err != nil {
		return err
	}
	return nil
}
