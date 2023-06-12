package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Validate[T any](payload T) error {

	validate := validator.New()

	if err := validate.Struct(payload); err != nil {

		validationErrors := err.(validator.ValidationErrors)

		var errorsBuilder strings.Builder

		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()
			errorMessage := fmt.Sprintf("Validation failed for field %s with tag '%s'\n", field, tag)
			errorsBuilder.WriteString(errorMessage)
		}

		errorsString := errorsBuilder.String()

		return &fiber.Error{Code: 400, Message: "Request Body Validation Failed: " + errorsString}

	}
	return nil
}
