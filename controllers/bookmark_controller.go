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

func GetPopulatedBookMarks(bookmarkType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")

		var bookmarks interface{}

		switch bookmarkType {
		case "post":
			var postBookmarks []models.PostBookmark
			if err := initializers.DB.
				Preload("PostItems.Post").
				Preload("PostItems.Post.User").
				Preload("PostItems.Post.RePost").
				Preload("PostItems.Post.RePost.User").
				Preload("PostItems.Post.TaggedUsers").
				Where("user_id = ?", loggedInUserID).
				Find(&postBookmarks).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			bookmarks = postBookmarks
		case "project":
			var projectBookmarks []models.ProjectBookmark
			if err := initializers.DB.
				Preload("ProjectItems.Project").
				Where("user_id = ?", loggedInUserID).
				Find(&projectBookmarks).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			var filteredBookmarks []models.ProjectBookmark

			for _, bookmark := range projectBookmarks {
				var projectItems []models.ProjectBookmarkItem
				for _, item := range bookmark.ProjectItems {
					if !item.Project.IsPrivate {
						projectItems = append(projectItems, item)
					}
				}
				bookmark.ProjectItems = projectItems
				filteredBookmarks = append(filteredBookmarks, bookmark)
			}

			bookmarks = filteredBookmarks
		case "opening":
			var openingBookmarks []models.OpeningBookmark
			if err := initializers.DB.
				Preload("OpeningItems.Opening").
				Preload("OpeningItems.Opening.Project").
				Where("user_id = ?", loggedInUserID).
				Find(&openingBookmarks).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			bookmarks = openingBookmarks
		default:
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid bookmarkType",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"status":    "success",
			"message":   "",
			"bookmarks": bookmarks,
		})
	}
}

func AddBookMark(bookmarkType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

		var ReqBody struct {
			Title string `json:"title"`
		}
		if err := c.BodyParser(&ReqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		var bookmark interface{}

		switch bookmarkType {
		case "post":
			postBookmark := models.PostBookmark{
				UserID: parsedLoggedInUserID,
				Title:  ReqBody.Title,
			}
			result := initializers.DB.Create(&postBookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
			bookmark = postBookmark
		case "project":
			projectBookmark := models.ProjectBookmark{
				UserID: parsedLoggedInUserID,
				Title:  ReqBody.Title,
			}
			result := initializers.DB.Create(&projectBookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
			bookmark = projectBookmark
		case "opening":
			openingBookmark := models.OpeningBookmark{
				UserID: parsedLoggedInUserID,
				Title:  ReqBody.Title,
			}
			result := initializers.DB.Create(&openingBookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
			bookmark = openingBookmark
		default:
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid bookmarkType",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"status":   "success",
			"message":  "",
			"bookmark": bookmark,
		})
	}
}

func UpdateBookMark(bookmarkType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookmarkID := c.Params("bookmarkID")
		loggedInUserID := c.GetRespHeader("loggedInUserID")

		var ReqBody struct {
			Title string `json:"title"`
		}
		if err := c.BodyParser(&ReqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		switch bookmarkType {
		case "post":
			var bookmark models.PostBookmark
			if err := initializers.DB.First(&bookmark, "id=? AND user_id=?", bookmarkID, loggedInUserID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
			}

			bookmark.Title = ReqBody.Title

			result := initializers.DB.Save(&bookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		case "project":
			var bookmark models.ProjectBookmark
			if err := initializers.DB.First(&bookmark, "id=? AND user_id=?", bookmarkID, loggedInUserID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
			}

			bookmark.Title = ReqBody.Title

			result := initializers.DB.Save(&bookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		case "opening":
			var bookmark models.OpeningBookmark
			if err := initializers.DB.First(&bookmark, "id=? AND user_id=?", bookmarkID, loggedInUserID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
			}

			bookmark.Title = ReqBody.Title

			result := initializers.DB.Save(&bookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		default:
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid bookmarkType",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Bookmark Edited.",
		})
	}
}

func DeleteBookMark(bookmarkType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookmarkID := c.Params("bookmarkID")
		loggedInUserID := c.GetRespHeader("loggedInUserID")

		switch bookmarkType {
		case "post":
			var bookmark models.PostBookmark
			if err := initializers.DB.First(&bookmark, "id=? AND user_id=?", bookmarkID, loggedInUserID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
			}

			result := initializers.DB.Delete(&bookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		case "project":
			var bookmark models.ProjectBookmark
			if err := initializers.DB.First(&bookmark, "id=? AND user_id=?", bookmarkID, loggedInUserID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
			}

			result := initializers.DB.Delete(&bookmark)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		case "opening":
			var bookmark models.OpeningBookmark
			if err := initializers.DB.First(&bookmark, "id=? AND user_id=?", bookmarkID, loggedInUserID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "No Bookmark of this ID found."}
			}

			result := initializers.DB.Delete(&bookmark)
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
