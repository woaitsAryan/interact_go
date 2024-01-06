package explore_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetRecommendedPosts(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	recommendations, err := utils.MLReq(loggedInUserID, config.POST_RECOMMENDATION)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, c.Path())
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"posts":  nil,
		})
	}

	var posts []models.Post

	if err := initializers.DB.
		Preload("User").
		Where("id IN ?", recommendations).
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.IncrementPostImpression(posts)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"posts":  posts,
	})
}

func GetRecommendedOpenings(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	recommendations, err := utils.MLReq(loggedInUserID, config.OPENING_RECOMMENDATION)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, c.Path())
		return c.Status(200).JSON(fiber.Map{
			"status":   "success",
			"openings": nil,
		})
	}

	var openings []models.Opening

	if err := initializers.DB.
		Preload("User").
		Where("id IN ?", recommendations).
		Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var filteredOpenings []models.Opening
	for _, opening := range openings {
		if opening.Project.UserID != parsedLoggedInUserID && !opening.Project.IsPrivate {
			filteredOpenings = append(filteredOpenings, opening)
		}
	}

	go routines.IncrementOpeningImpression(filteredOpenings)

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": filteredOpenings,
	})
}

func GetRecommendedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	recommendations, err := utils.MLReq(loggedInUserID, config.PROJECT_RECOMMENDATION)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, c.Path())
		return c.Status(200).JSON(fiber.Map{
			"status":   "success",
			"projects": nil,
		})
	}

	var projects []models.Project

	if err := initializers.DB.
		Preload("User").
		Preload("Memberships").
		Where("id IN ?", recommendations).
		Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.IncrementProjectImpression(projects)

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetRecommendedUsers(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.
		Preload("Profile").
		Where("active=? AND onboarding_completed=?", true, true).
		Where("organization_status=?", false).
		Where("verified=?", true).
		Where("username != email").
		Omit("phone_no").
		Omit("email").
		Select("*, (0.6 * no_followers - 0.4 * no_following + 0.3 * total_no_views) / (1 + EXTRACT(EPOCH FROM age(NOW(), created_at)) / 3600 / 24 / 21) AS weighted_average"). //! 21 days
		Order("weighted_average DESC, created_at ASC").
		Where("id <> ? AND organization_status = ?", loggedInUserID, false).
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.IncrementUserImpression(users)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetRecommendedEvents(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	recommendations, err := utils.MLReq(loggedInUserID, config.EVENT_RECOMMENDATION)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, c.Path())
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"events": nil,
		})
	}

	var events []models.Event

	if err := initializers.DB.
		Preload("Organization").
		Preload("Organization.User").
		Where("id IN ?", recommendations).
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.IncrementEventImpression(events)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"events": events,
	})
}

func GetRecommendedOrganizationalUsers(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	type UserWithOrganization struct {
		models.User
		Organization models.Organization `json:"organization"`
	}

	var users []models.User
	if err := searchedDB.
		Preload("Profile").
		Joins("LEFT JOIN organizations ON users.id = organizations.user_id").
		Where("users.active=?", true).
		Where("users.organization_status=?", true).
		Where("users.verified=?", true).
		Where("users.id <> ?", loggedInUserID).
		Where("users.username != users.email").
		Omit("users.phone_no").
		Omit("users.email").
		Select(`
        users.*,
        organizations.number_of_members,
        (0.6 * users.no_followers + 0.3 * users.total_no_views + 0.2 * organizations.number_of_members) / (1 + EXTRACT(EPOCH FROM age(NOW(), users.created_at)) / 3600 / 24 / 21) AS weighted_average
    `).
		Order("weighted_average DESC, users.created_at ASC").
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var usersWithOrganization []UserWithOrganization

	for _, user := range users {
		var organization models.Organization
		if err := initializers.DB.
			Where("user_id = ?", user.ID).
			First(&organization).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		userWithOrganization := UserWithOrganization{
			user,
			organization,
		}

		usersWithOrganization = append(usersWithOrganization, userWithOrganization)
	}

	go routines.IncrementUserImpression(users)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  usersWithOrganization,
	})
}
