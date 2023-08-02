package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AcceptApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("Opening").Preload("Opening.Project").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if application.Opening.UserID.String() != loggedInUserID {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	if application.Status == -1 {
		return &fiber.Error{Code: 400, Message: "Application is already Rejected."}
	}

	application.Status = 2
	result := initializers.DB.Save(&application)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.CreateMembershipAndSendNotification(&application)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Accepted.",
	})
}

func RejectApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("Opening").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if application.Opening.UserID.String() != loggedInUserID {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	if application.Status == 2 {
		return &fiber.Error{Code: 400, Message: "Application is already Accepted."}
	}

	application.Status = -1
	result := initializers.DB.Save(&application)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	notification := models.Notification{
		NotificationType: 7,
		UserID:           application.UserID,
		OpeningID:        &application.OpeningID,
	}

	if err := initializers.DB.Create(&notification).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Rejected.",
	})
}

func SetApplicationUnderReview(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("Opening").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if application.Opening.UserID.String() != loggedInUserID {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	if application.Status != 0 {
		return &fiber.Error{Code: 400, Message: "Cannot Set Under Review Now."}
	}

	application.Status = 1
	result := initializers.DB.Save(&application)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Under Review.",
	})
}
