package controllers

import (
	"strings"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func GetTrendingSearches(c *fiber.Ctx) error {
	var trendingSearches []string
	timeWindow := time.Now().Add(-24 * time.Hour)

	// Count the frequency of each normalized search query within the time window
	var searchCounts []struct {
		Query  string
		Counts int
	}
	initializers.DB.Table("search_queries").
		Select("LOWER(query) as query, COUNT(*) as counts"). // Ensure lowercase comparison
		Where("timestamp > ?", timeWindow).
		Group("LOWER(query)"). // Ensure grouping with lowercase
		Order("counts DESC").
		Limit(10). // You can adjust the number of trending searches you want to display
		Scan(&searchCounts)

	// Extract the search queries from the results
	for _, searchCount := range searchCounts {
		trendingSearches = append(trendingSearches, searchCount.Query)
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"searches": trendingSearches,
	})
}

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

func GetTrendingPosts(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var posts []models.Post

	searchedDB := API.Search(c, 2)(paginatedDB)

	if err := searchedDB.
		Preload("User").
		Select("*, (2 * no_likes + no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"posts":  posts,
	})
}

func GetTrendingOpenings(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var openings []models.Opening

	searchedDB := API.Search(c, 3)(paginatedDB)

	if err := searchedDB.Preload("Project").Order("created_at DESC").Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": openings,
	})
}

func GetProjectOpenings(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	var openings []models.Opening
	if err := initializers.DB.Preload("Project").Where("project_id = ?", projectID).Find(&openings).Order("created_at DESC").Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": openings,
	})
}

func GetTrendingProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.
		Select("*, (2 * no_likes + no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetRecommendedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Where("user_id <> ?", loggedInUserID).Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetMostLikedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Order("no_likes DESC").Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetRecentlyAddedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Order("created_at DESC").Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetLastViewedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var projectViewed []models.LastViewed

	paginatedDB := API.Paginator(c)(initializers.DB)
	if err := paginatedDB.Order("timestamp DESC").Preload("Project").Where("user_id=?", loggedInUserID).Find(&projectViewed).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var projects []models.Project

	for _, projectView := range projectViewed {
		projects = append(projects, projectView.Project)
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
		Select("*, (2 * no_followers - no_following) AS weighted_average").
		Order("weighted_average DESC").
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetRecommendedUsers(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.Order("created_at DESC").Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetSimilarUsers(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var user models.User
	initializers.DB.
		First(&user, "id = ?", userID)

	var users []models.User
	if err := initializers.DB.
		Where("id <> ?", userID).
		Where("tags && ?", pq.StringArray(user.Tags)).Limit(10).Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetSimilarProjects(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	var project models.Project
	initializers.DB.
		First(&project, "id = ?", projectID)

	var projects []models.Project
	if err := initializers.DB.
		Where("id <> ?", projectID).
		Where("category LIKE ?", "%"+project.Category+"%").
		Where("tags && ?", pq.StringArray(project.Tags)).Limit(10).Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}