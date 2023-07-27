package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetFeed(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var followings []models.FollowFollower
	if err := initializers.DB.Model(&models.FollowFollower{}).Where("follower_id = ?", loggedInUserID).Find(&followings).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	followingIDs := make([]uuid.UUID, len(followings))
	for i, following := range followings {
		followingIDs[i] = following.FollowedID
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.Preload("User").Where("user_id = ? OR user_id IN (?)", loggedInUserID, followingIDs).Order("created_at DESC").Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"feed":   posts,
	})
}
