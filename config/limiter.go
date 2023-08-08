package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RATE_LIMITER() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        1000,          // Maximum number of requests allowed within the duration
		Expiration: 1 * time.Hour, // Duration for which the limit applies
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address to differentiate clients
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			// What to do when the rate limit is reached
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests, please try again later.",
			})
		},
	})
}

const BODY_LIMIT = 5 * 1024 * 1024 // 5 MB
