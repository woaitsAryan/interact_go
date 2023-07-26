package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPostBookmark(c *fiber.Ctx) error {

	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var reqBody struct {
		Title string
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
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "New bookmark created.",
	})
}

func EditPostBookmark(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var bookmark models.PostBookmark
	if err := initializers.DB.First(&bookmark, "id = ? AND user_id=?", parsedBookmarkID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	var reqBody struct {
		Title string
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if reqBody.Title != "" {
		bookmark.Title = reqBody.Title
	}

	result := initializers.DB.Save(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark updated.",
	})
}

func DeletePostBookmark(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var bookmark models.PostBookmark
	if err := initializers.DB.First(&bookmark, "id = ? AND user_id=?", parsedBookmarkID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	result := initializers.DB.Delete(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark deleted.",
	})
}

func AddPostBookmarkItem(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var reqBody struct {
		PostID string `json:"postID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if reqBody.PostID != "" {
		parsedPostID, err := uuid.Parse(reqBody.PostID)
		if err != nil {
			return &fiber.Error{Code: 500, Message: "Invalid Post ID."}
		}

		item := models.PostBookmarkItem{
			PostBookmarkID: parsedBookmarkID,
			PostID:         parsedPostID,
		}

		result := initializers.DB.Create(&item)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Post added to the bookmark.",
		})
	}
	return &fiber.Error{Code: 400, Message: "Invalid Post ID."}
}

func RemovePostBookmarkItem(c *fiber.Ctx) error {
	bookmarkItemID := c.Params("bookmarkItemID")
	parsedBookmarkItemID, err := uuid.Parse(bookmarkItemID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var bookmarkItem models.PostBookmarkItem
	if err := initializers.DB.First(&bookmarkItem, "id = ?", parsedBookmarkItemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	result := initializers.DB.Delete(&bookmarkItem)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Post removed from the bookmark.",
	})
}
