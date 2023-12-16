package explore_controllers

import (
	"strconv"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func GetSimilarUsers(c *fiber.Ctx) error { //TODO ML Implementation
	username := c.Params("username")

	var user models.User
	initializers.DB.
		First(&user, "username = ?", username)

	var similarUsers []models.User
	if err := initializers.DB.
		Preload("Profile").
		Where("active=?", true).
		Where("organization_status=?", false).
		Where("id <> ?", username).
		Where("tags && ?", pq.StringArray(user.Tags)).
		Where("verified=?", true).
		Where("username != email").
		Omit("phone_no").
		Omit("email").
		Order("no_followers DESC").
		Limit(10).
		Find(&similarUsers).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.IncrementUserImpression(similarUsers)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  similarUsers,
	})
}

func GetSimilarProjects(c *fiber.Ctx) error {
	slug := c.Params("slug")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var project models.Project
	if err := initializers.DB.First(&project, "slug = ?", slug).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Project with this ID found."}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	recommendations, err := utils.MLReq(project.ID.String(), config.PROJECT_SIMILAR, limit, page)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, c.Path())
	}

	var projects []models.Project

	if len(recommendations) == 0 {
		if err := paginatedDB.
			Preload("User").
			Preload("Memberships").
			Where("id <> ?", project.ID).
			Where("is_private=?", false).
			Where("category = ? OR tags && ?", project.Category, pq.StringArray(project.Tags)).
			Order("total_no_views DESC").
			Find(&projects).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := initializers.DB.
			Preload("User").
			Preload("Memberships").
			Where("id IN ?", recommendations).
			Find(&projects).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	go routines.IncrementProjectImpression(projects)

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetSimilarEvents(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var event models.Event
	if err := initializers.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Event with this ID found."}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	recommendations, err := utils.MLReq(event.ID.String(), config.EVENT_SIMILAR, limit, page)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, c.Path())
	}

	var events []models.Event

	if len(recommendations) == 0 {
		if err := paginatedDB.
			Preload("Organization").
			Preload("Organization.User").
			Where("id <> ?", event.ID).
			Where("category = ? OR tags && ?", event.Category, pq.StringArray(event.Tags)).
			Order("no_views DESC").
			Find(&events).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := initializers.DB.
			Preload("Organization").
			Preload("Organization.User").
			Where("id IN ?", recommendations).
			Find(&events).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	go routines.IncrementEventImpression(events)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"events": events,
	})
}
