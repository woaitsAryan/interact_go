package cache

import (
	"context"
	"encoding/json"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.TODO()

func GetPost(key string) (*models.Post, error) {
	data, err := initializers.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "Error while getting post in cache."}
	}

	post := models.Post{}
	err = json.Unmarshal([]byte(data), &post)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "Error while unMarshalling post."}
	}
	return &post, nil
}

func SetPost(key string, post *models.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error while setting post in cache."}
	}

	initializers.RedisClient.Set(ctx, key, data, initializers.CacheExpirationTime)
	return nil
}
