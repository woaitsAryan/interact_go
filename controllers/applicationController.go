package controllers

import (
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

func GetApplication(c *fiber.Ctx) error { //! Only user and project applied members can get
	applicationID := c.Params("applicationID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("User").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "",
		"application": application,
	})
}

func GetAllApplicationsOfOpening(c *fiber.Ctx) error { //! Only project members can get
	openingID := c.Params("openingID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var applications []models.Application
	if err := initializers.DB.Preload("User").Where("opening_id=?", parsedOpeningID).Find(&applications).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"applications": applications,
	})
}

func AddApplication(c *fiber.Ctx) error {
	openingID := c.Params("openingID")
	userID := c.GetRespHeader("loggedInUserID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	parsedUserID, _ := uuid.Parse(userID)

	var application models.Application
	if err := initializers.DB.Where("opening_id=? AND user_id=?", parsedOpeningID, parsedUserID).First(&application).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "You already have applied for this opening."}
	}

	var reqBody schemas.ApplicationCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.ApplicationCreateSchema](reqBody); err != nil {
		return err
	}

	resumePath, err := utils.SaveFile(c, "resume", "project/openings/applications", false, 0, 0)
	if err != nil {
		return err
	}

	newApplication := models.Application{
		OpeningID: parsedOpeningID,
		UserID:    parsedUserID,
		Content:   reqBody.Content,
		Resume:    resumePath,
		Links:     reqBody.Links,
	}

	result := initializers.DB.Create(&newApplication)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the application."}
	}

	go routines.IncrementOpeningApplicationsAndSendNotification(parsedOpeningID, newApplication.ID, parsedUserID)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "New Application Added",
	})
}

func DeleteApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.First(&application, "user_id=? AND id = ?", loggedInUserID, parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&application)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the application."}
	}

	parsedOpeningID := application.OpeningID
	go routines.DecrementOpeningApplications(parsedOpeningID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Deleted",
	})
}
