package helpers

import (
	"fmt"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (err AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Error: %v", err.Code, err.Message, err.Err)
}

func ErrorHandler(c *fiber.Ctx, err error) error {

	Code := 500
	Message := config.SERVER_ERROR
	Error := err

	if e, ok := err.(*AppError); ok {
		Code = e.Code
		Message = e.Message
		Error = e.Err
	}

	if Message == config.DATABASE_ERROR {
		LogDatabaseError("Database Error", Error, c.Path())
	} else if Code == 500 {
		LogServerError("Server Error", Error, c.Path())
	}

	return c.Status(Code).JSON(fiber.Map{
		"status":  "failed",
		"message": Message,
	})
}
