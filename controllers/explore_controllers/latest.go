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
	"gorm.io/gorm"
)

func GetLatestPosts(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var posts []models.Post

	searchedDB := API.Search(c, 2)(paginatedDB)

	if err := searchedDB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("RePost").
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
		Select("*, posts.id, posts.created_at").
		Order("posts.created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementPostImpression(posts)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"posts":  posts,
	})
}

func GetLatestProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("Memberships").
		Order("created_at DESC").
		Where("is_private = ?", false).
		Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementProjectImpression(projects)

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}
