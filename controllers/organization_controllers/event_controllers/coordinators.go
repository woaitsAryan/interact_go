package event_controllers

import (
	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddEventCoordinators(c *fiber.Ctx) error {
	eventID := c.Params("eventID")
	parsedOrgID, _ := uuid.Parse(c.Params("orgID"))
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	var reqBody struct {
		UserIDs []string `json:"userIDs"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var event models.Event
	if err := initializers.DB.First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var users []models.User

	for _, userID := range reqBody.UserIDs {
		var membership models.OrganizationMembership
		if err := initializers.DB.Preload("User").First(&membership, "user_id = ? AND organization_id=?", userID, event.OrganizationID).Error; err != nil {
			continue
		}

		users = append(users, membership.User)
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "AddEventCoordinators")
		}
	}()

	if err := tx.Model(&event).Association("Coordinators").Clear(); err != nil {
		return err
	}

	event.Coordinators = users

	if err := tx.Save(&event).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 16, nil, nil, &event.ID, nil, nil, nil, nil, nil, nil, nil, "")
	go cache.RemoveEvent(event.ID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Coordinators added",
	})
}

func RemoveEventCoordinators(c *fiber.Ctx) error {
	eventID := c.Params("eventID")
	parsedOrgID, _ := uuid.Parse(c.Params("orgID"))
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	var event models.Event
	if err := initializers.DB.First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Model(&event).Association("Coordinators").Clear(); err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 17, nil, nil, &event.ID, nil, nil, nil, nil, nil, nil, nil, "")
	go cache.RemoveEvent(event.ID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Coordinators removed",
	})
}
