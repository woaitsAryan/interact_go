package project_controllers

import (
	"github.com/Pratham-Mishra04/interact/cache"
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
	if err := initializers.DB.Preload("Project").First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, err := uuid.Parse(loggedInUserID)

	if err == nil && parsedLoggedInUserID != opening.UserID && parsedLoggedInUserID != opening.Project.UserID {
		go routines.UpdateLastViewedOpening(parsedLoggedInUserID, opening.ID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"opening": opening,
	})
}

func GetAllOpeningsOfProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	var project models.Project
	if err := initializers.DB.Where("id=?", projectID).First(&project).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Cannot perform this action"}
	}

	var openings []models.Opening
	if err := initializers.DB.Where("project_id=?", projectID).Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"openings": openings,
	})
}

func AddOpening(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	userID := c.GetRespHeader("loggedInUserID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? AND user_id=?", parsedProjectID, parsedUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.OpeningCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	if err := helpers.Validate[schemas.OpeningCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	newOpening := models.Opening{
		ProjectID:   parsedProjectID,
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Tags:        reqBody.Tags,
		UserID:      parsedUserID,
	}

	result := initializers.DB.Create(&newOpening)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	projectMemberID := c.GetRespHeader("projectMemberID")
	parsedID, _ := uuid.Parse(projectMemberID)
	go routines.MarkProjectHistory(project.ID, parsedID, 3, nil, &newOpening.ID, nil, nil, nil, "")

	go cache.RemoveProject(project.Slug)
	go cache.RemoveProject("-workspace--" + project.Slug)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "New Opening Added",
		"opening": newOpening,
	})
}

func EditOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")
	userID := c.GetRespHeader("loggedInUserID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	parsedUserID, err := uuid.Parse(userID)
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
	if err := initializers.DB.Preload("Project").First(&opening, "id = ? AND user_id=?", parsedOpeningID, parsedUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			for _, application := range pendingApplications {
				application.Status = -1
				result := initializers.DB.Save(&application)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
				}
			}
		}
	}

	result := initializers.DB.Save(&opening)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	projectMemberID := c.GetRespHeader("projectMemberID")
	parsedID, _ := uuid.Parse(projectMemberID)
	go routines.MarkProjectHistory(opening.ProjectID, parsedID, 4, nil, &opening.ID, nil, nil, nil, "")

	go cache.RemoveProject(opening.Project.Slug)
	go cache.RemoveProject("-workspace--" + opening.Project.Slug)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Opening Updated",
	})
}

func DeleteOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")
	userID := c.GetRespHeader("loggedInUserID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var opening models.Opening
	if err := initializers.DB.Preload("Project").First(&opening, "id = ? AND user_id=?", parsedOpeningID, parsedUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	projectID := opening.ProjectID
	projectSlug := opening.Project.Slug

	result := initializers.DB.Delete(&opening)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	projectMemberID := c.GetRespHeader("projectMemberID")
	parsedID, _ := uuid.Parse(projectMemberID)
	go routines.MarkProjectHistory(projectID, parsedID, 5, nil, nil, nil, nil, nil, opening.Title)

	go cache.RemoveProject(projectSlug)
	go cache.RemoveProject("-workspace--" + projectSlug)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Opening Deleted",
	})
}
