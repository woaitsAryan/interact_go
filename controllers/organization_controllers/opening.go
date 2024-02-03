package organization_controllers

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

func GetOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var opening models.Opening
	if err := initializers.DB.Preload("Organization").First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, err := uuid.Parse(loggedInUserID)

	if err == nil && parsedLoggedInUserID != opening.UserID && parsedLoggedInUserID != opening.Organization.UserID {
		go routines.UpdateLastViewedOpening(parsedLoggedInUserID, opening.ID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"opening": opening,
	})
}

func GetAllOpeningsOfOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var openings []models.Opening
	if err := initializers.DB.Where("organization_id=?", orgID).Order("created_at DESC").Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"openings": openings,
	})
}

func AddOpening(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	orgMemberID := c.GetRespHeader("orgMemberID")

	parsedOrgMemberID, err := uuid.Parse(orgMemberID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}
	parsedOrgID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "id = ?", parsedOrgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody schemas.OpeningCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	if err := helpers.Validate[schemas.OpeningCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	newOpening := models.Opening{
		OrganizationID: &parsedOrgID,
		Title:          reqBody.Title,
		Description:    reqBody.Description,
		Tags:           reqBody.Tags,
		UserID:         parsedOrgMemberID,
	}

	result := initializers.DB.Create(&newOpening)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 24, nil, nil, nil, nil, nil, nil, nil, &newOpening.ID, "")

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "New Opening Added",
		"opening": newOpening,
	})
}

func EditOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")
	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")
	parsedOrgID, _ := uuid.Parse(orgID)

	parsedOrgMemberID, err := uuid.Parse(orgMemberID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var reqBody schemas.OpeningEditSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.OpeningEditSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	var opening models.Opening
	if err := initializers.DB.Preload("Organization").First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.Description != "" {
		opening.Description = reqBody.Description
	}
	if reqBody.Tags != nil {
		opening.Tags = reqBody.Tags
	}
	if reqBody.Active != nil {
		opening.Active = *reqBody.Active
		if !opening.Active {
			var pendingApplications []models.Application
			if err := initializers.DB.Find(&pendingApplications, "opening_id AND (status=0 OR status=1)", parsedOpeningID).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			for _, application := range pendingApplications {
				application.Status = -1
				result := initializers.DB.Save(&application)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
			}
		}
	}

	result := initializers.DB.Save(&opening)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 26, nil, nil, nil, nil, nil, nil, nil, &opening.ID, "")

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Opening Updated",
	})
}

func DeleteOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")
	orgID := c.Params("orgID")
	parsedOrgID, _ := uuid.Parse(orgID)
	orgMemberID := c.GetRespHeader("orgMemberID")
	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var opening models.Opening
	if err := initializers.DB.Preload("Organization").First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	result := initializers.DB.Delete(&opening)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 25, nil, nil, nil, nil, nil, nil, nil, nil, opening.Title)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Opening Deleted",
	})
}
