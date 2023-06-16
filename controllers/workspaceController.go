package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
)

func GetMyProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	searchedDB := API.Search(c, 1)(initializers.DB)

	var projects []models.Project
	if err := searchedDB.Where("user_id = ?", loggedInUserID).Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projects,
	})
}

func GetMyContributingProjects(c *fiber.Ctx) error { //!Add search here
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var memberships []models.Membership
	if err := initializers.DB.Preload("Project").Select("project_id").Where("user_id = ?", loggedInUserID).Find(&memberships).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

func GetMyApplications(c *fiber.Ctx) error { //! Add search here
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var applications []models.Application
	if err := initializers.DB.Where("user_id=?", loggedInUserID).Find(&applications).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"applications": applications,
	})
}
