package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func LikePost(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	postID := c.Params("postID")
	parsedPostID, err := uuid.Parse(postID)

	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var like models.UserPostLike
	err = initializers.DB.Where("user_id=? AND post_id=?", parsedLoggedInUserID, parsedPostID).First(&like).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserPostLike{
				PostID: parsedPostID,
				UserID: parsedLoggedInUserID,
			}

			result := initializers.DB.Create(&likeModel)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			go routines.IncrementPostLikesAndSendNotification(parsedPostID, parsedLoggedInUserID)

		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
		go routines.DecrementPostLikes(parsedPostID)

	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Liked/Unliked.",
	})
}

func LikeProject(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID, _ := uuid.Parse(loggedInUserID)

	projectID := c.Params("projectID")
	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var like models.UserProjectLike
	if err := initializers.DB.Where("user_id=? AND project_id=?", userID, parsedProjectID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserProjectLike{
				ProjectID: parsedProjectID,
				UserID:    userID,
			}

			result := initializers.DB.Create(&likeModel)

			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			go routines.IncrementProjectLikesAndSendNotification(parsedProjectID, userID)
		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
		go routines.DecrementProjectLikes(parsedProjectID)
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project Liked/Unliked.",
	})
}

func LikeComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID, _ := uuid.Parse(loggedInUserID)

	commentID := c.Params("commentID")
	parsedCommentID, err := uuid.Parse(commentID)

	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var comment models.Comment
	if err := initializers.DB.Where("id = ?", parsedCommentID).First(&comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var like models.UserCommentLike
	if err := initializers.DB.Where("user_id=? AND comment_id=?", userID, parsedCommentID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserCommentLike{
				PostID:    comment.PostID,
				ProjectID: comment.ProjectID,
				CommentID: comment.ID,
				UserID:    userID,
			}
			result := initializers.DB.Create(&likeModel)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			go routines.IncrementCommentLikes(parsedCommentID, userID)
		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
		go routines.DecrementCommentLikes(parsedCommentID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Liked/Unliked.",
	})
}
