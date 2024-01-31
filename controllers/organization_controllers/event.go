package organization_controllers

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
		Preload("Organization.User").
		Preload("Coordinators").
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
		CoverPic:       picName,
		Description:    reqBody.Description,
		Tags:           reqBody.Tags,
		Category:       reqBody.Category,
		Links:          reqBody.Links,
		OrganizationID: parsedOrgID,
		StartTime:      startTime,
		EndTime:        endTime,
		Location:       reqBody.Location,
	}

	result := initializers.DB.Create(&event)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 0, nil, nil, &event.ID, nil, nil, nil, nil, nil, "")
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

	var reqBody schemas.EventUpdateSchema
	c.BodyParser(&reqBody)

	picName, err := utils.UploadImage(c, "coverPic", helpers.EventClient, 1920, 1080)
	if err != nil {
		return err
	}
	oldEventPic := event.CoverPic

	if reqBody.Tagline != "" {
		event.Tagline = reqBody.Tagline
	}
	if picName != "" {
		event.CoverPic = picName
	}
	if reqBody.Category != "" {
		event.Category = reqBody.Category
	}
	if reqBody.Description != "" {
		event.Description = reqBody.Description
	}
	if reqBody.Location != "" {
		event.Location = reqBody.Location
	}
	if reqBody.Tags != nil {
		event.Tags = reqBody.Tags
	}
	if reqBody.Links != nil {
		event.Links = reqBody.Links
	}
	if reqBody.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, reqBody.StartTime)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Start Time."}
		}
		event.StartTime = startTime
	}
	if reqBody.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, reqBody.EndTime)
		if err != nil || endTime.Before(event.StartTime) {
			return &fiber.Error{Code: 400, Message: "Invalid End Time."}
		}

		event.EndTime = endTime
	}

	if err := initializers.DB.Save(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.EventClient, oldEventPic)
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 2, nil, nil, &event.ID, nil, nil, nil, nil, nil,  "")
	routines.GetImageBlurHash(c, "coverPic", &event)
	go cache.RemoveEvent(event.ID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Event updated successfully",
		"event":   event,
	})
}

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

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 16, nil, nil, &event.ID, nil, nil, nil, nil, nil, "")
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

	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 17, nil, nil, &event.ID, nil, nil, nil, nil, nil, "")
	go cache.RemoveEvent(event.ID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Coordinators removed",
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

	var event models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	eventPic := event.CoverPic

	if err := initializers.DB.Delete(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DeleteFromBucket(helpers.EventClient, eventPic)
	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 1, nil, nil, nil, nil, nil, nil, nil, nil, event.Title)
	go routines.DecrementOrgEvent(parsedOrgID)
	go cache.RemoveEvent(event.ID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}

func AddOtherOrg(c *fiber.Ctx) error{
	var reqBody schemas.CoHostEventSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}
	eventID := c.Params("eventID")
	// orgID := c.Params("orgID")
	// parsedOrgID, _ := uuid.Parse(orgID)

	var event models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	parsedCoOwnOrgId , err := uuid.Parse(reqBody.OrganizationID)
	if err != nil{	
		return &fiber.Error{Code: 400, Message: "Invalid Co own organization ID."}
	}

	var CoOwnOrganization models.Organization
	if err := initializers.DB.Where("id = ?", reqBody.OrganizationID).First(&CoOwnOrganization).Error; err != nil{
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
		}
		return helpers.AppError{Code:500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	event.CoOwnedBy = append(event.CoOwnedBy, CoOwnOrganization)

	var invitation models.Invitation
	invitation.OrganizationID = &parsedCoOwnOrgId
	invitation.UserID = CoOwnOrganization.UserID
	invitation.Title = "Your organization has been invited to cohost an event!"
	invitation.EventID = &event.ID
	result := initializers.DB.Create(&invitation)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	if err := initializers.DB.Save(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization added as a cohost",
		"event":   event,
	})
}

func RemoveOtherOrg(c *fiber.Ctx) error{
	var reqBody schemas.CoHostEventSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	eventID := c.Params("eventID")

	var event models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	parsedCoOwnOrgId , err := uuid.Parse(reqBody.OrganizationID)
	if err != nil{	
		return &fiber.Error{Code: 400, Message: "Invalid Co own organization ID."}
	}

	var CoOwnOrganization models.Organization
	if err := initializers.DB.Where("id = ?", reqBody.OrganizationID).First(&CoOwnOrganization).Error; err != nil{
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
		}
		return helpers.AppError{Code:500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}	
	found := false
	for i, coOwner := range event.CoOwnedBy {
		if(coOwner.ID == parsedCoOwnOrgId){
			event.CoOwnedBy = append(event.CoOwnedBy[:i], event.CoOwnedBy[i+1:]...)
			found = true
		}
	}
	if(!found){
		return &fiber.Error{Code: 400, Message: "No co own organization of this ID found"}
	}

	if err := initializers.DB.Save(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DecrementOrgEvent(parsedCoOwnOrgId)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization removed as a cohost",
		"event":   event,
	})
}
