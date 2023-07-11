package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddProjectBookmark(c *fiber.Ctx) error {

	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var reqBody struct {
		Title string
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
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating bookmark."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "New bookmark created.",
	})
}

func EditProjectBookmark(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var bookmark models.ProjectBookmark
	if err := initializers.DB.First(&bookmark, "id = ? AND user_id=?", parsedBookmarkID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
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
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating bookmark."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark updated.",
	})
}

func DeleteProjectBookmark(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var bookmark models.ProjectBookmark
	if err := initializers.DB.First(&bookmark, "id = ? AND user_id=?", parsedBookmarkID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&bookmark)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting bookmark."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark deleted.",
	})
}

func AddProjectBookmarkItem(c *fiber.Ctx) error {
	bookmarkID := c.Params("bookmarkID")
	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var reqBody struct {
		ProjectID string `json:"projectID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if reqBody.ProjectID != "" {
		parsedProjectID, err := uuid.Parse(reqBody.ProjectID)
		if err != nil {
			return &fiber.Error{Code: 500, Message: "Invalid Project ID."}
		}

		item := models.ProjectBookmarkItem{
			ProjectBookmarkID: parsedBookmarkID,
			ProjectID:         parsedProjectID,
		}

		result := initializers.DB.Create(&item)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while creating bookmark item."}
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Project added to the bookmark.",
		})
	}
	return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
}

func RemoveProjectBookmarkItem(c *fiber.Ctx) error {
	bookmarkItemID := c.Params("bookmarkItemID")
	parsedBookmarkItemID, err := uuid.Parse(bookmarkItemID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid bookmark ID."}
	}

	var bookmarkItem models.ProjectBookmarkItem
	if err := initializers.DB.First(&bookmarkItem, "id = ?", parsedBookmarkItemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Bookmark Item of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&bookmarkItem)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting bookmark item."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Project removed from the bookmark.",
	})
}
