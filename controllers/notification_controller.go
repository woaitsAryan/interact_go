package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetNotifications(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var notifications []models.Notification
	if err := paginatedDB.
		Preload("User").
		Preload("Sender").
		Preload("Post").
		Preload("Project").
		Preload("Opening").
		Preload("Application").
		Where("user_id=?", loggedInUserID).
		Order("created_at DESC").
		Find(&notifications).
		Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"message":       "",
		"notifications": notifications,
	})
}

func DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("notificationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var notification models.Notification
	if err := initializers.DB.First(&notification, "id = ? AND user_id=?", notificationID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Notification of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	result := initializers.DB.Delete(&notification)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Notification Deleted",
	})
}
