package organization_controllers

import (
	"time"

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

	var event models.Event
	if err := initializers.DB.
		Preload("Organization").
		Where("id = ?", eventID).
		First(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

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

	picName, err := utils.UploadFile(c, "coverPic", helpers.EventClient, 2560, 2560)
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	// orgMemberID := c.GetRespHeader("orgMemberID")
	// if orgMemberID != "" {
	// 	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
	// 	go routines.MarkOrganizationHistory(event.ID, parsedOrgMemberID, -1, nil, nil, nil, nil, nil)
	// }
	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 0, nil, nil, &event.ID, nil, nil)

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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.EventUpdateSchema
	c.BodyParser(&reqBody)

	// picName, err := utils.SaveFile(c, "coverPic", "project/coverPics", true, 2560, 2560)
	picName, err := utils.UploadFile(c, "coverPic", helpers.EventClient, 2560, 2560)
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.EventClient, oldEventPic)
	}

	// orgMemberID := c.GetRespHeader("orgMemberID")
	// if orgMemberID != "" {
	// 	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
	// 	go routines.MarkOrganizationHistory(project.ID, parsedOrgMemberID, 2, nil, nil, nil, nil, nil)
	// }
	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 2, nil, nil, &event.ID, nil, nil)

	//TODO setup event cache
	// cache.RemoveProject(project.Slug)
	// cache.RemoveProject("-workspace--" + project.Slug)

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

	var event models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	eventPic := event.CoverPic

	if err := initializers.DB.Delete(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.DeleteFromBucket(helpers.EventClient, eventPic)

	// orgMemberID := c.GetRespHeader("orgMemberID")
	// if orgMemberID != "" {
	// 	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
	// 	go routines.MarkOrganizationHistory(project.ID, parsedOrgMemberID, 2, nil, nil, nil, nil, nil)
	// }
	go routines.MarkOrganizationHistory(parsedOrgID, parsedUserID, 1, nil, nil, &event.ID, nil, nil)

	//TODO setup event cache
	// cache.RemoveProject(project.Slug)
	// cache.RemoveProject("-workspace--" + project.Slug)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}
