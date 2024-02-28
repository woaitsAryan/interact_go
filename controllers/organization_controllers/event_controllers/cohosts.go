package event_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetEventCoHosts(c *fiber.Ctx) error {
	eventID := c.Params("eventID")
	orgID := c.Params("orgID")

	var event models.Event
	if err := initializers.DB.Preload("CoOwnedBy").
		Preload("CoOwnedBy.User").
		Where("id = ? AND organization_id=?", eventID, orgID).
		First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var invitations []models.Invitation
	if err := initializers.DB.Preload("User").
		Where("event_id=?", eventID).
		Find(&invitations).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "",
		"coHosts":     event.CoOwnedBy,
		"invitations": invitations,
	})
}

func AddCoHostOrgs(c *fiber.Ctx) error {
	var reqBody schemas.CoHostEventSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}
	eventID := c.Params("eventID")
	orgID := c.Params("orgID")

	var event models.Event
	if err := initializers.DB.Preload("CoOwnedBy").Where("id = ? AND organization_id=?", eventID, orgID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	coOwnedOrgIDs := []string{}
	for _, coHostOrg := range event.CoOwnedBy {
		coOwnedOrgIDs = append(coOwnedOrgIDs, coHostOrg.ID.String())
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "AddCoHostOrg")
		}
	}()

	for _, userID := range reqBody.UserIDs {
		var CoOwnOrganization models.Organization
		if err := tx.Where("user_id = ?", userID).First(&CoOwnOrganization).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		if CoOwnOrganization.ID == event.OrganizationID || utils.Contains(coOwnedOrgIDs, orgID) {
			continue
		}

		var existingInvitation models.Invitation
		if err := initializers.DB.First(&existingInvitation, "event_id=? AND user_id=? AND status=0", event.ID, CoOwnOrganization.UserID).Error; err == nil {
			continue
		}

		invitation := models.Invitation{
			UserID:  CoOwnOrganization.UserID,
			Title:   event.Title,
			EventID: &event.ID,
		}

		result := tx.Create(&invitation)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitations Sent!",
	})
}

func RemoveCoHostOrgs(c *fiber.Ctx) error {
	var reqBody schemas.CoHostEventSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	eventID := c.Params("eventID")
	orgID := c.Params("orgID")

	var event models.Event
	if err := initializers.DB.Preload("CoOwnedBy").Where("id = ? AND organization_id=?", eventID, orgID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	toRemoveOrgIDs := []uuid.UUID{}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "AddCoHostOrg")
		}
	}()

	for _, userID := range reqBody.UserIDs {
		var CoOwnOrganization models.Organization
		if err := initializers.DB.Where("user_id = ?", userID).First(&CoOwnOrganization).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		toRemoveOrgIDs = append(toRemoveOrgIDs, CoOwnOrganization.ID)

		if err := tx.Model(&event).Association("CoOwnedBy").Delete(&CoOwnOrganization); err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, toRemoveOrgID := range toRemoveOrgIDs {
		go routines.DecrementOrgEvent(toRemoveOrgID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Organizations removed as from cohosts",
	})
}

func LeaveCoHostOrg(c *fiber.Ctx) error {
	eventID := c.Params("eventID")
	orgID := c.Params("orgID")

	var event models.Event
	if err := initializers.DB.Preload("CoOwnedBy").Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	parsedCoOwnOrgId, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Co own organization ID."}
	}

	var organization models.Organization
	if err := initializers.DB.Where("id = ?", parsedCoOwnOrgId).First(&organization).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	found := false
	for i, coOwner := range event.CoOwnedBy {
		if coOwner.ID == organization.ID {
			event.CoOwnedBy = append(event.CoOwnedBy[:i], event.CoOwnedBy[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
	}

	if err := initializers.DB.Save(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DecrementOrgEvent(parsedCoOwnOrgId)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization left as a cohost",
	})
}
