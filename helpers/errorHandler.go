package helpers

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
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

	// Send custom error page
	// err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
	// if err != nil {
	//     return ctx.Status(500).SendString(config.SERVER_ERROR)
	// }

	return c.Status(code).JSON(fiber.Map{
		"status":  "failed",
		"message": message,
	})

}

func CatchError(statusCode int, message string, err error) *fiber.Error {
	// Log the error to the console
	if message == config.DATABASE_ERROR || message == config.SERVER_ERROR {
		initializers.Logger.Errorw(message, "Message", err.Error(), "Error ", err)
	}

	return &fiber.Error{
		Code:    statusCode,
		Message: message,
	}
}
