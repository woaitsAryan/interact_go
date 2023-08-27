package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddBookMarkItem(bookmarkType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookmarkID := c.Params("bookmarkID")
		parsedBookmarkID, err := uuid.Parse(bookmarkID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID"}
		}

		var reqBody struct {
			ItemID string `json:"itemID"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		parsedItemID, err := uuid.Parse(reqBody.ItemID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID"}
		}

		var bookmarkItem interface{}

		switch bookmarkType {
		case "post":
			postBookmarkItem := models.PostBookmarkItem{
				PostBookmarkID: parsedBookmarkID,
				PostID:         parsedItemID,
			}

			result := initializers.DB.Create(&postBookmarkItem)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			bookmarkItem = postBookmarkItem
		case "project":
			projectBookmarkItem := models.ProjectBookmarkItem{
				ProjectBookmarkID: parsedBookmarkID,
				ProjectID:         parsedItemID,
			}

			result := initializers.DB.Create(&projectBookmarkItem)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			bookmarkItem = projectBookmarkItem
		case "openings":
			openingBookmarkItem := models.OpeningBookmarkItem{
				OpeningBookmarkID: parsedBookmarkID,
				OpeningID:         parsedItemID,
			}

			result := initializers.DB.Create(&openingBookmarkItem)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			bookmarkItem = openingBookmarkItem
		default:
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid bookmarkType",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"status":       "success",
			"message":      "Bookmark Item Created.",
			"bookmarkItem": bookmarkItem,
		})
	}
}

func DeleteBookMarkItem(bookmarkType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookmarkItemID := c.Params("bookmarkItemID")

		switch bookmarkType {
		case "post":
			var bookmarkItem models.PostBookmarkItem
			err := initializers.DB.First(&bookmarkItem, "id=?", bookmarkItemID).Error
			if err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
			}

			result := initializers.DB.Delete(&bookmarkItem)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		case "project":
			var bookmarkItem models.ProjectBookmarkItem
			err := initializers.DB.First(&bookmarkItem, "id=?", bookmarkItemID).Error
			if err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
			}

			result := initializers.DB.Delete(&bookmarkItem)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		case "openings":
			var bookmarkItem models.OpeningBookmarkItem
			err := initializers.DB.First(&bookmarkItem, "id=?", bookmarkItemID).Error
			if err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
			}

			result := initializers.DB.Delete(&bookmarkItem)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		default:
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid bookmarkType",
			})
		}

		return c.Status(204).JSON(fiber.Map{
			"status":  "success",
			"message": "Bookmark Deleted.",
		})
	}
}
