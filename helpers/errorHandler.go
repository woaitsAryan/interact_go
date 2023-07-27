package helpers

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	code := 500
	message := config.SERVER_ERROR

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	if message == config.DATABASE_ERROR {
		LogDatabaseError("Database Error", err, c.Path())
	}

	if code == 500 {
		LogServerError("Server Error", err, c.Path())
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "failed",
		"message": message,
	})

}
