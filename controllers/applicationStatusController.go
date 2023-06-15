package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AcceptApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("Opening").Preload("Opening.Project").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if application.Status == -1 {
		return &fiber.Error{Code: 400, Message: "Application is already Rejected."}
	}

	result := initializers.DB.Save(&application)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the application."}
	}

	membership := models.Membership{
		ProjectID: application.Opening.ProjectID,
		UserID:    application.UserID,
		Role:      "",
		Title:     application.Opening.Title,
	}

	result = initializers.DB.Create(&membership)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating membership."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Accepted.",
	})
}

func RejectApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if application.Status == 2 {
		return &fiber.Error{Code: 400, Message: "Application is already Accepted."}
	}

	result := initializers.DB.Save(&application)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the application."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Rejected.",
	})
}

func SetApplicationUnderReview(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if application.Status != 0 {
		return &fiber.Error{Code: 400, Message: "Cannot Set Under Review Now."}
	}

	result := initializers.DB.Save(&application)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the application."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Under Review.",
	})
}
