package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetNotifications(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID, _ := uuid.Parse(loggedInUserID)

	paginatedDB := API.Paginator(c)(initializers.DB)

	var notifications []models.Notification
	if err := paginatedDB.
		Preload("User").
		Preload("Sender").
		Preload("Post").
		Preload("Project").
		Preload("Opening").
		Preload("Application").
		Where("user_id=?", userID).
		Find(&notifications).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"message":       "",
		"notifications": notifications,
	})
}

func DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("notificationID")

	parsedNotificationID, err := uuid.Parse(notificationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var notification models.Notification
	if err := initializers.DB.First(&notification, "id = ?", parsedNotificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Notification of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&notification)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the notification."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Notification Deleted",
	})
}
