package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
)

func GetTrendingPosts(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var posts []models.Post

	postUserSelectedDB := utils.PostSelectConfig(paginatedDB.Preload("User"))

	if err := postUserSelectedDB.Order("created_at DESC").Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Posts."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"posts":  posts,
	})
}

func GetTrendingOpenings(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var openings []models.Opening

	if err := paginatedDB.Preload("Project").Order("created_at DESC").Find(&openings).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Openings."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": openings,
	})
}
