package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetBookMarks(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var postBookmarks []models.PostBookmark
	err := initializers.DB.Preload("PostItems").Find(&postBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error."}
	}

	var projectBookmarks []models.ProjectBookmark
	err = initializers.DB.Preload("ProjectItems").Find(&projectBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":           "success",
		"message":          "",
		"postBookmarks":    postBookmarks,
		"projectBookmarks": projectBookmarks,
	})
}

func AddPostBookMark(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	bookmark := models.PostBookmark{
		UserID: parsedUserID,
		Title:  reqBody.Title,
	}

	result := initializers.DB.Create(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the bookmark."}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "success",
		"message":  "Bookmark Created.",
		"bookmark": bookmark,
	})
}

func AddProjectBookMark(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	bookmark := models.ProjectBookmark{
		UserID: parsedUserID,
		Title:  reqBody.Title,
	}

	result := initializers.DB.Create(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the bookmark."}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "success",
		"message":  "Bookmark Created.",
		"bookmark": bookmark,
	})
}

func DeletePostBookMark(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.PostBookmark
	err = initializers.DB.First(&bookmark, "id=?", parsedBookmarkID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
	}

	result := initializers.DB.Delete(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the bookmark."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}

func DeleteProjectBookMark(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.PostBookmark
	err = initializers.DB.First(&bookmark, "id=?", parsedBookmarkID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
	}

	result := initializers.DB.Delete(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the bookmark."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}
