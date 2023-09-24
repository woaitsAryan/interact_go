package controllers

import (
	"strings"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
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
		Preload("RePost").
		Preload("RePost.User").
		Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
		Select("*, posts.id, posts.created_at, (2 * no_likes + no_comments + 5 * no_shares) AS weighted_average").
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
	//TODO add openings of the matched project name
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	paginatedDB := API.Paginator(c)(initializers.DB)
	var openings []models.Opening

	searchedDB := API.Search(c, 3)(paginatedDB)

	if err := searchedDB.Preload("Project").
		Order("no_of_applications DESC").
		Where("active=true").
		Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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

func GetRecommendedOpenings(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	paginatedDB := API.Paginator(c)(initializers.DB)
	var openings []models.Opening

	if err := paginatedDB.Preload("Project").
		Where("active=true").
		Order("created_at DESC").Find(&openings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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

func GetTrendingProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.
		Preload("User").
		Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Where("is_private = ?", false).
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

	if err := searchedDB.
		Preload("User").
		Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Where("user_id <> ? AND is_private = ?", loggedInUserID, false).
		Find(&projects).Error; err != nil {
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

	if err := searchedDB.
		Preload("User").
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

func GetRecentlyAddedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.
		Preload("User").
		Order("created_at DESC").
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

	var projectViewed []models.LastViewed

	paginatedDB := API.Paginator(c)(initializers.DB)
	if err := paginatedDB.
		Preload("User").
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

func GetTrendingUsers(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.
		Where("active=?", true).
		Where("organization_status=?", false).
		Omit("phone_no").
		Omit("email").
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
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.Where("active=?", true).
		Omit("phone_no").
		Omit("email").
		Order("no_followers DESC").
		Where("id <> ? AND organization_status = ?", loggedInUserID, false).
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetSimilarUsers(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	initializers.DB.
		First(&user, "username = ?", username)

	var similarUsers []models.User
	if err := initializers.DB.
		Where("active=?", true).
		Where("organization_status=?", false).
		Where("id <> ?", username).
		Where("tags && ?", pq.StringArray(user.Tags)).
		Omit("phone_no").
		Omit("email").
		Order("no_followers DESC").
		Limit(10).
		Find(&similarUsers).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  similarUsers,
	})
}

func GetSimilarProjects(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var project models.Project
	initializers.DB.
		First(&project, "slug = ?", slug)

	var projects []models.Project
	if err := initializers.DB.
		Where("id <> ?", project.ID).
		Where("category LIKE ?", "%"+project.Category+"%").
		Where("tags && ?", pq.StringArray(project.Tags)).
		Where("is_private=?", false).
		Limit(10).
		Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}
