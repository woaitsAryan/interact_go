package organization_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.First(organization, "id=?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"organization": organization,
	})
}

func GetOrganizationTasks(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.
		Preload("Memberships").
		Preload("Memberships.User").
		Find(&organization, "id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var tasks []models.Task
	if err := initializers.DB.
		Preload("Users").
		Preload("SubTasks").
		Preload("SubTasks.Users").
		Find(&tasks, "organization_id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"tasks":        tasks,
		"organization": organization,
	})
}
