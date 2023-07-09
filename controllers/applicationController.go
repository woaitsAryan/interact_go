package controllers

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetApplication(c *fiber.Ctx) error {
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

func GetAllApplicationsOfOpening(c *fiber.Ctx) error {
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

	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", parsedOpeningID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
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

	opening.NoOfApplications++
	result = initializers.DB.Save(&opening)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the opening."}
	}

	// notification := models.Notification{
	// 	NotificationType:       5,
	// 	ReceiverID:             opening.UserID,
	// 	SenderID:               parsedUserID,
	// 	ConcernedOpeningID:     opening.ID,
	// 	ConcernedApplicationID: newApplication.ID,
	// 	ConcernedProjectID:     opening.ProjectID,
	// }

	// if err := initializers.DB.Create(&notification).Error; err != nil {
	// 	return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
	// }

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "New Application Added",
	})
}

func DeleteApplication(c *fiber.Ctx) error {
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

	parsedOpeningID := application.OpeningID

	result := initializers.DB.Delete(&application)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the application."}
	}

	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", parsedOpeningID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
	}

	opening.NoOfApplications--
	result = initializers.DB.Save(&opening)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the opening."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Deleted",
	})
}
