package utils

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Exists(modelDoc *gorm.Model, id uuid.UUID) (*gorm.Model, error) {
	if err := initializers.DB.First(modelDoc, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return modelDoc, &fiber.Error{Code: 400, Message: "No Document of this ID found."}
		}
		return modelDoc, helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return modelDoc, nil
}
