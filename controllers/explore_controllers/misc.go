package explore_controllers

import (
	"strings"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddSearchQuery(c *fiber.Ctx) error {
	var reqBody struct {
		Search string `json:"search"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}
	searchQuery := models.SearchQuery{
		Query: strings.ToLower(strings.TrimSpace(reqBody.Search)),
	}
	initializers.DB.Create(&searchQuery)
	return c.Status(201).JSON(fiber.Map{
		"status": "success",
	})
}

func GetProjectOpenings(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var project models.Project
	if err := initializers.DB.Where("slug = ? AND is_private = ?", slug, false).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(400).JSON(fiber.Map{
				"status":  "success",
				"message": "Project Not Found",
			})
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var openings []models.Opening
	if err := initializers.DB.
		Preload("Project").
		Where("project_id = ? AND active=true", project.ID).
		Order("created_at DESC").
		Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": openings,
	})
}

func GetOrgEvents(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var events []models.Event
	if err := paginatedDB.
		Preload("Organization").
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"events":  events,
	})
}

func GetMostLikedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.
		Preload("User").
		Preload("Memberships").
		Order("no_likes DESC").
		Where("is_private = ?", false).
		Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetLastViewedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var projectViewed []models.LastViewedProjects

	paginatedDB := API.Paginator(c)(initializers.DB)
	if err := paginatedDB.
		Preload("User").
		Preload("Memberships").
		Order("timestamp DESC").
		Preload("Project").
		Where("user_id=?", loggedInUserID).
		Find(&projectViewed).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var projects []models.Project

	for _, projectView := range projectViewed {
		if !projectView.Project.IsPrivate {
			projects = append(projects, projectView.Project)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetOrganizationalUser(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	if err := initializers.DB.First(&user, "username=?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "user_id=?", user.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"user":         user,
		"organization": organization,
	})
}
