package project_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMyProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	searchedDB := API.Search(c, 1)(initializers.DB)

	var projects []models.Project
	if err := searchedDB.Where("user_id = ?", loggedInUserID).Order("created_at DESC").Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projects,
	})
}

func GetMyContributingProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var memberships []models.Membership
	if err := initializers.DB.Preload("Project").Preload("Project.User", func(db *gorm.DB) *gorm.DB {
		return db.Select(select_fields.User)
	}).Where("user_id = ?", loggedInUserID).Order("created_at DESC").Find(&memberships).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var projects []models.Project
	for _, membership := range memberships {
		projects = append(projects, membership.Project)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projects,
	})
}

func GetMyApplications(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var applications []models.Application
	if err := initializers.DB.
		Preload("Opening").
		Preload("Project", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.Project)
		}).
		Preload("Organization").
		Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Where("user_id=?", loggedInUserID).
		Order("created_at DESC").
		Find(&applications).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"applications": applications,
	})
}

func GetMyMemberships(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var memberships []models.Membership
	if err := initializers.DB.Where("user_id = ?", loggedInUserID).Find(&memberships).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "",
		"memberships": memberships,
	})
}
