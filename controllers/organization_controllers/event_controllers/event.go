package event_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	eventInCache, err := cache.GetEvent(eventID)
	if err == nil {
		go routines.UpdateEventViews(eventInCache.ID)
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "",
			"event":   eventInCache,
		})
	}

	var event models.Event
	if err := initializers.DB.
		Preload("Organization").
		Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("Coordinators", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("CoOwnedBy").
		Preload("CoOwnedBy.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Where("id = ?", eventID).
		First(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.UpdateEventViews(event.ID)
	go cache.SetEvent(event.ID.String(), &event)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"event":   event,
	})
}

func GetEventHistory(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	var history []models.EventHistory
	if err := initializers.DB.
		Preload("User").
		Where("event_id = ?", eventID).
		First(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"history": history,
	})
}

func GetPopulatedOrgEvents(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var events []models.Event
	if err := paginatedDB.
		Preload("Organization").
		Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("CoOwnedBy").
		Preload("CoOwnedBy.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("Coordinators", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Joins("LEFT JOIN co_owned_events ON co_owned_events.event_id = events.id").
		Where("events.organization_id = ? OR co_owned_events.organization_id = ?", orgID, orgID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"events":  events,
	})
}
func AddEvent(c *fiber.Ctx) error {
	var reqBody schemas.EventCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.EventCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	picName, err := utils.UploadImage(c, "coverPic", helpers.EventClient, 1920, 1080)
	if err != nil {
		return err
	}

	startTime, err := time.Parse(time.RFC3339, reqBody.StartTime)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Start Time."}
	}

	endTime, err := time.Parse(time.RFC3339, reqBody.EndTime)
	if err != nil || endTime.Before(startTime) {
		return &fiber.Error{Code: 400, Message: "Invalid End Time."}
	}

	event := models.Event{
		Title:          reqBody.Title,
		Tagline:        reqBody.Tagline,
		Description:    reqBody.Description,
		Tags:           reqBody.Tags,
		Category:       reqBody.Category,
		Links:          reqBody.Links,
		OrganizationID: parsedOrgID,
		StartTime:      startTime,
		EndTime:        endTime,
		Location:       reqBody.Location,
	}

	if picName != "" {
		event.CoverPic = picName
	}

	result := initializers.DB.Create(&event)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 0, nil, nil, &event.ID, nil, nil, nil, nil, nil, nil, nil, "")
	go routines.IncrementOrgEvent(parsedOrgID)
	routines.GetImageBlurHash(c, "coverPic", &event)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Event Added",
		"event":   event,
	})
}

func UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	var event models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	history := models.EventHistory{
		EventID: event.ID,
		UserID:  parsedUserID,
	}

	var reqBody schemas.EventUpdateSchema
	c.BodyParser(&reqBody)

	picName, err := utils.UploadImage(c, "coverPic", helpers.EventClient, 1920, 1080)
	if err != nil {
		return err
	}
	oldEventPic := event.CoverPic

	if reqBody.Tagline != "" {
		event.Tagline = reqBody.Tagline
		history.Tagline = true
	}
	if picName != "" {
		event.CoverPic = picName
		history.CoverPic = true
	}
	if reqBody.Category != "" {
		event.Category = reqBody.Category
		history.Category = true
	}
	if reqBody.Description != "" {
		event.Description = reqBody.Description
		history.Description = true
	}
	if reqBody.Location != "" {
		event.Location = reqBody.Location
		history.Location = true
	}
	if reqBody.Tags != nil {
		event.Tags = reqBody.Tags
		history.Tags = true
	}
	if reqBody.Links != nil {
		event.Links = reqBody.Links
		history.Links = true
	}
	if reqBody.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, reqBody.StartTime)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Start Time."}
		}
		event.StartTime = startTime
		history.StartTime = true
	}
	if reqBody.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, reqBody.EndTime)
		if err != nil || endTime.Before(event.StartTime) {
			return &fiber.Error{Code: 400, Message: "Invalid End Time."}
		}

		event.EndTime = endTime
		history.EndTime = true
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Save(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Create(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.EventClient, oldEventPic)
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 2, nil, nil, &event.ID, nil, nil, nil, nil, nil, nil, nil, "")
	routines.GetImageBlurHash(c, "coverPic", &event)
	go cache.RemoveEvent(event.ID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Event updated successfully",
		"event":   event,
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventID")
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var event models.Event
	if err := initializers.DB.Preload("CoOwnedBy").Preload("Invitations").Where("id = ? AND organization_id=?", eventID, parsedOrgID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	eventPic := event.CoverPic
	coHosts := event.CoOwnedBy

	if err := tx.Model(&event).Association("CoOwnedBy").Clear(); err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Model(&event).Association("Coordinators").Clear(); err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, invitation := range event.Invitations {
		if err := tx.Delete(&invitation).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	var history models.EventHistory
	if err := tx.Where("event_id = ?", eventID).First(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Delete(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Delete(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DeleteFromBucket(helpers.EventClient, eventPic)
	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, event.Title)
	go cache.RemoveEvent(event.ID.String())

	go routines.DecrementOrgEvent(parsedOrgID)
	for _, org := range coHosts {
		go routines.DecrementOrgEvent(org.ID)
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}
