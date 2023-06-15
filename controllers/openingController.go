package controllers

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
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
	if err := initializers.DB.First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var applications []models.Application

	initializers.DB.Where("opening_id=?", parsedOpeningID).Find(&applications)

	opening.NoOfApplications = len(applications)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"opening": opening,
	})
}

func GetAllOpeningsOfProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var openings []models.Opening
	if err := initializers.DB.Where("project_id=?", parsedProjectID).Find(&openings).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

	var reqBody schemas.OpeningCreateScheam
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.OpeningCreateScheam](reqBody); err != nil {
		return err
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
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the opening."}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "New Opening Added",
	})
}

func EditOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var reqBody schemas.OpeningEditScheam
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.OpeningEditScheam](reqBody); err != nil {
		return err
	}

	var opening models.Opening
	if err := initializers.DB.First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if reqBody.Description != "" {
		opening.Description = reqBody.Description
	}
	if reqBody.Tags != nil {
		opening.Tags = reqBody.Tags
	}

	result := initializers.DB.Save(&opening)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the opening."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Opening Updated",
	})
}

func DeleteOpening(c *fiber.Ctx) error {
	openingID := c.Params("openingID")

	parsedOpeningID, err := uuid.Parse(openingID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var opening models.Opening
	if err := initializers.DB.First(&opening, "id = ?", parsedOpeningID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Opening of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&opening)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the opening."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Opening Deleted",
	})
}
