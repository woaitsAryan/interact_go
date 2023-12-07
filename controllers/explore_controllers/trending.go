package explore_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetTrendingSearches(c *fiber.Ctx) error {
	var trendingSearches []string
	// timeWindow := time.Now().Add(-10 * 24 * time.Hour)

	var searchCounts []struct {
		Query  string
		Counts int
	}
	searchedDB := API.Search(c, 4)(initializers.DB)
	searchedDB.Table("search_queries").
		Select("LOWER(query) as query, COUNT(*) as counts"). // Ensure lowercase comparison
		// Where("timestamp > ?", timeWindow).
		Group("LOWER(query)"). // Ensure grouping with lowercase
		Order("counts DESC").
		Limit(15).
		Scan(&searchCounts)

	for _, searchCount := range searchCounts {
		trendingSearches = append(trendingSearches, searchCount.Query)
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"searches": trendingSearches,
	})
}

func GetTrendingPosts(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)
	var posts []models.Post

	searchedDB := API.Search(c, 2)(paginatedDB)

	if loggedInUserID == "" {
		if err := searchedDB.
			Preload("User").
			Preload("RePost").
			Preload("RePost.User").
			Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
			Select("*, posts.id, posts.created_at, (2 * no_likes + no_comments + 5 * no_shares) / (1 + EXTRACT(EPOCH FROM age(NOW(), posts.created_at)) / 3600 / 24 / 7) AS weighted_average"). //! 7 days
			Order("weighted_average DESC, posts.created_at ASC").
			Find(&posts).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := searchedDB.
			Preload("User").
			Preload("RePost").
			Preload("RePost.User").
			Where("user_id <> ?", loggedInUserID).
			Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
			Select("*, posts.id, posts.created_at, (2 * no_likes + no_comments + 5 * no_shares) / (1 + EXTRACT(EPOCH FROM age(NOW(), posts.created_at)) / 3600 / 24 / 7) AS weighted_average"). //! 7 days
			Order("weighted_average DESC, posts.created_at ASC").
			Find(&posts).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"posts":  posts,
	})
}

func GetTrendingOpenings(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	paginatedDB := API.Paginator(c)(initializers.DB)
	var openings []models.Opening

	searchStr := c.Query("search", "")
	if searchStr == "" {
		if err := paginatedDB.Preload("Project").
			Joins("JOIN projects ON openings.project_id = projects.id").
			Where("openings.active=true").
			Select("openings.*, (projects.total_no_views * 0.5 + openings.no_of_applications * 0.3) / (1 + EXTRACT(EPOCH FROM age(NOW(), openings.created_at)) / 3600 / 24 / 15) AS t_ratio").
			Order("t_ratio DESC").
			Find(&openings).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		searchedDB := API.Search(c, 3)(paginatedDB)

		if err := searchedDB.Preload("Project").
			Where("openings.active=true").
			Select("openings.*, (projects.total_no_views * 0.5 + openings.no_of_applications * 0.3) / (1 + EXTRACT(EPOCH FROM age(NOW(), openings.created_at)) / 3600 / 24 / 15) AS t_ratio").
			Order("t_ratio DESC").
			Find(&openings).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}
	var filteredOpenings []models.Opening
	for _, opening := range openings {
		if opening.Project.UserID != parsedLoggedInUserID && !opening.Project.IsPrivate {
			filteredOpenings = append(filteredOpenings, opening)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": filteredOpenings,
	})
}

func GetTrendingProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if loggedInUserID != "" {
		if err := searchedDB.
			Preload("User").
			Preload("Memberships").
			Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
			Order("weighted_average DESC").
			Where("user_id <> ? AND is_private = ?", loggedInUserID, false).
			Find(&projects).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := searchedDB.
			Preload("User").
			Preload("Memberships").
			Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
			Order("weighted_average DESC").
			Where("is_private = ?", false).
			Find(&projects).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetTrendingUsers(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.
		Where("active=?", true).
		Where("organization_status=?", false).
		Where("verified=?", true).
		Where("username != email").
		Omit("phone_no").
		Omit("email").
		Select("*, (0.6 * no_followers - 0.4 * no_following + 0.3 * total_no_views) / (1 + EXTRACT(EPOCH FROM age(NOW(), created_at)) / 3600 / 24 / 21) AS weighted_average"). //! 21 days
		Order("weighted_average DESC, created_at ASC").
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetTrendingEvents(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var events []models.Event

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.
		Preload("Organization").
		Preload("Organization.User").
		Select("*, (no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"events": events,
	})
}
