package helpers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	code := 500
	message := "Internal server error."

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	// Send custom error page
	// err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
	// if err != nil {
	//     return ctx.Status(500).SendString("Internal Server Error")
	// }

	return c.Status(code).JSON(fiber.Map{
		"status":  "failed",
		"message": message,
	})

}
