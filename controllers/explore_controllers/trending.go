package explore_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetTrendingSearches(c *fiber.Ctx) error {
	var trendingSearches []string
	// timeWindow := time.Now().Add(-10 * 24 * time.Hour)

	var searchCounts []struct {
		Query  string
		Counts int
	}
	searchedDB := API.Search(c, 5)(initializers.DB)
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

	query := searchedDB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select(select_fields.User)
	}).Preload("RePost").
		Preload("RePost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("RePost.TaggedUsers", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ShorterUser)
		}).
		Preload("TaggedUsers", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ShorterUser)
		}).
		Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
		Where("is_flagged=?", false).
		Select("*, posts.id, posts.created_at, (2 * no_likes + no_comments + 5 * no_shares) / (1 + EXTRACT(EPOCH FROM age(NOW(), posts.created_at)) / 3600 / 24 / 7) AS weighted_average"). //* 7 days
		Order("weighted_average DESC, posts.created_at ASC")

	if loggedInUserID != "" {
		query = query.Where("user_id <> ?", loggedInUserID)
	}

	if err := query.Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementPostImpression(posts)

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

		filteredDB := API.Filter(c, 4)(paginatedDB)

		if err := filteredDB.
			Preload("Project", func(db *gorm.DB) *gorm.DB {
				return db.Select(select_fields.Project)
			}).
			Preload("Organization").
			Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
				return db.Select(select_fields.User)
			}).
			Where("active=true").
			Select("openings.*, (no_of_applications * 0.3) / (1 + EXTRACT(EPOCH FROM age(NOW(), openings.created_at)) / 3600 / 24 / 15) AS t_ratio").
			Order("t_ratio DESC").
			Find(&openings).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	} else {
		searchedDB := API.Search(c, 3)(paginatedDB)

		filteredDB := API.Filter(c, 4)(searchedDB)

		if err := filteredDB.
			Preload("Project", func(db *gorm.DB) *gorm.DB {
				return db.Select(select_fields.Project)
			}).
			Preload("Organization").
			Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
				return db.Select(select_fields.User)
			}).
			Where("active=true").
			Select("openings.*, (no_of_applications * 0.3) / (1 + EXTRACT(EPOCH FROM age(NOW(), openings.created_at)) / 3600 / 24 / 15) AS t_ratio").
			Order("t_ratio DESC").
			Find(&openings).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
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

func GetTrendingProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	filteredDB := API.Filter(c, 1)(searchedDB)

	if loggedInUserID != "" {
		if err := filteredDB.
			Preload("User").
			Preload("Memberships").
			Select("*, (projects.total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
			Order("weighted_average DESC").
			Where("user_id <> ? AND is_private = ?", loggedInUserID, false).
			Find(&projects).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	} else {
		if err := filteredDB.
			Preload("User").
			Preload("Memberships").
			Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
			Order("weighted_average DESC").
			Where("is_private = ?", false).
			Find(&projects).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	go routines.IncrementProjectImpression(projects)

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetTrendingUsers(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	searchedDB := API.Search(c, 0)(paginatedDB)
	filteredDB := API.Filter(c, 2)(searchedDB)

	var users []models.User
	if err := filteredDB.
		Preload("Profile").
		Where("active=? AND onboarding_completed=?", true, true).
		Where("verified=?", true).
		Where("username != email").
		Where("organization_status=?", false).
		Where("username != users.email").
		Omit("phone_no").
		Omit("users.email").
		Select("*, (0.6 * no_followers - 0.4 * no_following + 0.3 * total_no_views) / (1 + EXTRACT(EPOCH FROM age(NOW(), created_at)) / 3600 / 24 / 21) AS weighted_average"). //* 21 days
		Order("weighted_average DESC, created_at ASC").
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementUserImpression(users)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetTrendingEvents(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var events []models.Event

	searchedDB := API.Search(c, 4)(paginatedDB)

	filteredDB := API.Filter(c, 3)(searchedDB)

	if err := filteredDB.
		Preload("Organization").
		Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Select("*, events.id, (no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementEventImpression(events)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"events": events,
	})
}

func GetTrendingOrganizationalUsers(c *fiber.Ctx) error {
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
		Where("users.active=? AND users.onboarding_completed=?", true, true).
		Where("users.organization_status=?", true).
		Where("users.verified=?", true).
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var usersWithOrganization []UserWithOrganization

	for _, user := range users {
		var organization models.Organization
		if err := initializers.DB.
			Where("user_id = ?", user.ID).
			First(&organization).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
