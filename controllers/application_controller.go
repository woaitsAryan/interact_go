package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetApplication(applicationType string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		applicationID := c.Params("applicationID")
		loggedInUserID := c.GetRespHeader("loggedInUserID")

		parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

		parsedApplicationID, err := uuid.Parse(applicationID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID"}
		}

		var application models.Application

		switch applicationType {
		case "project":
			if err := initializers.DB.Preload("User").Preload("Project").Preload("Opening").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			if application.UserID != parsedLoggedInUserID && application.Project.UserID != parsedLoggedInUserID {
				return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
			}
		case "org":
			if err := initializers.DB.Preload("User").Preload("Organization").Preload("Opening").First(&application, "id = ?", parsedApplicationID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			if application.UserID != parsedLoggedInUserID && application.Organization.UserID != parsedLoggedInUserID {
				return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
			}
		}

		if application.IncludeEmail {
			application.Email = application.User.Email
		}
		if application.IncludeResume {
			application.Resume = application.User.Resume
		}

		return c.Status(200).JSON(fiber.Map{
			"status":      "success",
			"message":     "",
			"application": application,
		})
	}
}

func GetAllApplicationsOfOpening(c *fiber.Ctx) error {

	openingID := c.Params("openingID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var applications []models.Application
	if err := initializers.DB.Preload("User").Where("opening_id=?", parsedOpeningID).Order("created_at DESC").Find(&applications).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for i, application := range applications {
		if application.IncludeEmail {
			applications[i].Email = application.User.Email
		}
		if application.IncludeResume {
			applications[i].Resume = application.User.Resume
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"applications": applications,
	})

}

func AddApplication(applicationType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		openingID := c.Params("openingID")
		userID := c.GetRespHeader("loggedInUserID")

		parsedOpeningID, err := uuid.Parse(openingID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID"}
		}

		parsedUserID, _ := uuid.Parse(userID)

		var user models.User
		if err := initializers.DB.Where("id=?", parsedUserID).First(&user).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
		if !user.Verified {
			return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
		}
		if user.OrganizationStatus {
			return &fiber.Error{Code: 400, Message: "Organizations cannot perform this action."}
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
			return &fiber.Error{Code: 400, Message: err.Error()}
		}

		newApplication := models.Application{
			OpeningID:     parsedOpeningID,
			UserID:        parsedUserID,
			Content:       reqBody.Content,
			Links:         reqBody.Links,
			IncludeEmail:  reqBody.IncludeEmail,
			IncludeResume: reqBody.IncludeResume,
		}

		var opening models.Opening

		switch applicationType {
		case "project":
			if err := initializers.DB.Preload("Project").Where("id = ? AND active=true", parsedOpeningID).First(&opening).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
			}

			if opening.Project.UserID == parsedUserID {
				return &fiber.Error{Code: 400, Message: "You already are the creator of this project."}
			}

			var membership models.Membership
			if err := initializers.DB.Where("project_id=? AND user_id=?", opening.ProjectID, parsedUserID).First(&membership).Error; err == nil {
				return &fiber.Error{Code: 400, Message: "You already are a collaborator of this project."}
			}

			var invitation models.Invitation
			err = initializers.DB.Where("user_id=? AND project_id=? AND status=0", user.ID, opening.ProjectID).First(&invitation).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "You already are invited to join this project."}
			}

			newApplication.ProjectID = opening.ProjectID
		case "org":
			if err := initializers.DB.Preload("Organization").Where("id = ? AND active=true", parsedOpeningID).First(&opening).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
			}

			if opening.Organization.UserID == parsedUserID {
				return &fiber.Error{Code: 400, Message: "You already are the creator of this organization."}
			}

			var membership models.OrganizationMembership
			if err := initializers.DB.Where("organization_id=? AND user_id=?", opening.OrganizationID, parsedUserID).First(&membership).Error; err == nil {
				return &fiber.Error{Code: 400, Message: "You already are a member of this organization."}
			}

			var invitation models.Invitation
			err = initializers.DB.Where("user_id=? AND organization_id=? AND status=0", user.ID, opening.OrganizationID).First(&invitation).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "You already are invited to join this organization."}
			}

			newApplication.OrganizationID = opening.OrganizationID
		}

		result := initializers.DB.Create(&newApplication)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
		}

		go routines.IncrementOpeningApplicationsAndSendNotification(parsedOpeningID, newApplication.ID, parsedUserID)

		// go cache.RemoveProject("-workspace--" + opening.Project.Slug)

		return c.Status(201).JSON(fiber.Map{
			"status":  "success",
			"message": "New Application Added",
		})
	}
}

func DeleteApplication(c *fiber.Ctx) error {
	applicationID := c.Params("applicationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedApplicationID, err := uuid.Parse(applicationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var application models.Application
	if err := initializers.DB.Preload("Project").First(&application, "user_id=? AND id = ?", loggedInUserID, parsedApplicationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Application of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	result := initializers.DB.Delete(&application)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	parsedOpeningID := application.OpeningID
	go routines.DecrementOpeningApplications(parsedOpeningID)

	// go cache.RemoveProject("-workspace--" + opening.Project.Slug)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Application Deleted",
	})
}
