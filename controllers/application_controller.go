package controllers

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

func GetApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("User").Preload("Project").Preload("Opening").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if application.UserID != parsedLoggedInUserID && application.Project.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "",
		"application": application,
	})
}

func GetAllApplicationsOfOpening(c *fiber.Ctx) error { //! Save memberships in redux for frontend security
	openingID := c.Params("openingID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var applications []models.Application
	if err := initializers.DB.Preload("User").Where("opening_id=?", parsedOpeningID).Find(&applications).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var opening models.Opening
	if err := initializers.DB.Preload("Project").First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if len(applications) > 0 {
		var memberships []models.Membership
		if err := initializers.DB.Where("project_id=?", applications[0].ProjectID).Find(&memberships).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
		check := false

		for _, membership := range memberships {
			if membership.UserID == parsedLoggedInUserID {
				check = true
			}
		}

		if opening.Project.UserID == parsedLoggedInUserID {
			check = true
		}

		if !check {
			return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
		}
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

	var user models.User
	if err := initializers.DB.Where("id=?", parsedUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	if !user.Verified {
		return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
	}

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

	var opening models.Opening
	if err := initializers.DB.Where("id = ? AND active=true", parsedOpeningID).First(&opening).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
	}

	var membership models.Membership
	if err := initializers.DB.Where("project_id=? AND user_id=?", opening.ProjectID, parsedUserID).First(&membership).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "You already are a collaborator of this project."}
	}

	newApplication := models.Application{
		OpeningID: parsedOpeningID,
		ProjectID: opening.ProjectID,
		UserID:    parsedUserID,
		Content:   reqBody.Content,
		Resume:    resumePath,
		Links:     reqBody.Links,
	}

	result := initializers.DB.Create(&newApplication)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	result := initializers.DB.Delete(&application)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	parsedOpeningID := application.OpeningID
	go routines.DecrementOpeningApplications(parsedOpeningID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Deleted",
	})
}
