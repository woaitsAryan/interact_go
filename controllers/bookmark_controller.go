package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var projectBookmarks []models.ProjectBookmark
	err = initializers.DB.Preload("ProjectItems").Find(&projectBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var openingBookmarks []models.OpeningBookmark
	err = initializers.DB.Preload("OpeningItems").Find(&openingBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":           "success",
		"message":          "",
		"postBookmarks":    postBookmarks,
		"projectBookmarks": projectBookmarks,
		"openingBookmarks": openingBookmarks,
	})
}

func GetPopulatedPostBookMarks(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var postBookmarks []models.PostBookmark
	err := initializers.DB.Preload("PostItems.Post").Preload("PostItems.Post.User").Find(&postBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"bookmarks": postBookmarks,
	})
}

func GetPopulatedProjectBookMarks(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var projectBookmarks []models.ProjectBookmark
	err := initializers.DB.Preload("ProjectItems.Project").Find(&projectBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"bookmarks": projectBookmarks,
	})
}

func GetPopulatedOpeningBookMarks(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var openingBookmarks []models.OpeningBookmark
	err := initializers.DB.Preload("OpeningItems.Opening").Find(&openingBookmarks, "user_id=?", parsedUserID).Error
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"bookmarks": openingBookmarks,
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}
	return c.Status(201).JSON(fiber.Map{
		"status":   "success",
		"message":  "Bookmark Created.",
		"bookmark": bookmark,
	})
}

// TODO create a single add, delete bookmark controller which takes in the input for the type of bookmark
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "success",
		"message":  "Bookmark Created.",
		"bookmark": bookmark,
	})
}

func AddOpeningBookMark(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	bookmark := models.OpeningBookmark{
		UserID: parsedUserID,
		Title:  reqBody.Title,
	}

	result := initializers.DB.Create(&bookmark)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "success",
		"message":  "Bookmark Created.",
		"bookmark": bookmark,
	})
}

func DeletePostBookMark(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.PostBookmark
	err = initializers.DB.First(&bookmark, "id=? AND user_id=?", parsedBookmarkID, loggedInUserID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
	}

	result := initializers.DB.Delete(&bookmark)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}

func DeleteProjectBookMark(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.ProjectBookmark
	err = initializers.DB.First(&bookmark, "id=? AND user_id=?", parsedBookmarkID, loggedInUserID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
	}

	result := initializers.DB.Delete(&bookmark)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}

func DeleteOpeningBookMark(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.OpeningBookmark
	err = initializers.DB.First(&bookmark, "id=? AND user_id=?", parsedBookmarkID, loggedInUserID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
	}

	result := initializers.DB.Delete(&bookmark)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}
