package config

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     initializers.CONFIG.FRONTEND_URL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	})
}

func API_CHECKER(c *fiber.Ctx) error {
	REQ_TOKEN := c.Get("Authentication")
	if REQ_TOKEN != initializers.CONFIG.API_TOKEN {
		return &fiber.Error{Code: 403, Message: "Cannot access the API"}
	}
	return c.Next()
}
