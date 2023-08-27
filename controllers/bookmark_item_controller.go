package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPostBookMarkItem(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var reqBody struct {
		PostID string `json:"postID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedPostID, err := uuid.Parse(reqBody.PostID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.PostBookmark
	if err := initializers.DB.First(&bookmark, "id = ?", parsedBookmarkID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	bookmarkItem := models.PostBookmarkItem{
		PostBookmarkID: parsedBookmarkID,
		PostID:         parsedPostID,
	}

	result := initializers.DB.Create(&bookmarkItem)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Item Created.",
	})
}

func AddProjectBookMarkItem(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var reqBody struct {
		ProjectID string `json:"projectID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedProjectID, err := uuid.Parse(reqBody.ProjectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.ProjectBookmark
	if err := initializers.DB.First(&bookmark, "id = ?", parsedBookmarkID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	bookmarkItem := models.ProjectBookmarkItem{
		ProjectBookmarkID: parsedBookmarkID,
		ProjectID:         parsedProjectID,
	}

	result := initializers.DB.Create(&bookmarkItem)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Item Created.",
	})
}

func AddOpeningBookMarkItem(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var reqBody struct {
		OpeningID string `json:"openingID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedOpeningID, err := uuid.Parse(reqBody.OpeningID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmark models.OpeningBookmark
	if err := initializers.DB.First(&bookmark, "id = ?", parsedBookmarkID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	bookmarkItem := models.OpeningBookmarkItem{
		OpeningBookmarkID: parsedBookmarkID,
		OpeningID:         parsedOpeningID,
	}

	result := initializers.DB.Create(&bookmarkItem)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Item Created.",
	})
}

func DeletePostBookMarkItem(c *fiber.Ctx) error {
	bookmarkItemID := c.Params("bookmarkItemID")
	parsedBookmarkItemID, err := uuid.Parse(bookmarkItemID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmarkItem models.PostBookmarkItem
	err = initializers.DB.First(&bookmarkItem, "id=?", parsedBookmarkItemID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
	}

	result := initializers.DB.Delete(&bookmarkItem)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}

func DeleteProjectBookMarkItem(c *fiber.Ctx) error {
	projectBookmarkItemID := c.Params("bookmarkItemID")
	parsedBookmarkItemID, err := uuid.Parse(projectBookmarkItemID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmarkItem models.ProjectBookmarkItem
	err = initializers.DB.First(&bookmarkItem, "id=?", parsedBookmarkItemID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
	}

	result := initializers.DB.Delete(&bookmarkItem)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}

func DeleteOpeningBookMarkItem(c *fiber.Ctx) error {
	projectBookmarkItemID := c.Params("bookmarkItemID")
	parsedBookmarkItemID, err := uuid.Parse(projectBookmarkItemID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var bookmarkItem models.OpeningBookmarkItem
	err = initializers.DB.First(&bookmarkItem, "id=?", parsedBookmarkItemID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
	}

	result := initializers.DB.Delete(&bookmarkItem)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark Deleted.",
	})
}
