package controllers

import (
	"log"

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

	searchedDB := API.SearchPosts(c)(postUserSelectedDB)

	if err := searchedDB.Find(&posts).Order("created_at DESC").Error; err != nil {
		log.Fatal(err)
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

	searchedDB := API.Search(c, 3)(paginatedDB)

	if err := searchedDB.Preload("Project").Order("created_at DESC").Find(&openings).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Openings."}
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

	if err := searchedDB.Order("created_at DESC").Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
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
	if err := searchedDB.Order("created_at DESC").Find(&users).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}
