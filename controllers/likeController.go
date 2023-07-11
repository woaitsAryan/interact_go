package controllers

import (
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
				return &fiber.Error{Code: 500, Message: "Internal Server Error while adding the like."}
			}
			go routines.IncrementPostLikesAndSendNotification(parsedPostID, parsedLoggedInUserID)

		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
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
				return &fiber.Error{Code: 500, Message: "Internal Server Error while adding the like."}
			}
			go routines.IncrementProjectLikesAndSendNotification(parsedProjectID, userID)
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}
		go routines.DecrementProjectLikes(parsedProjectID)
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project Liked/Unliked.",
	})
}

func LikePostComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID, _ := uuid.Parse(loggedInUserID)

	commentID := c.Params("commentID")
	parsedCommentID, err := uuid.Parse(commentID)

	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var like models.UserPostCommentLike
	if err := initializers.DB.Where("user_id=? AND post_comment_id=?", userID, parsedCommentID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserPostCommentLike{
				PostCommentID: parsedCommentID,
				UserID:        userID,
			}
			result := initializers.DB.Create(&likeModel)
			if result.Error != nil {
				return &fiber.Error{Code: 500, Message: "Internal Server Error while adding the like."}
			}
			go routines.IncrementPostCommentLikes(parsedCommentID, userID)
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}
		go routines.DecrementPostCommentLikes(parsedCommentID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Liked/Unliked.",
	})
}

func LikeProjectComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID, _ := uuid.Parse(loggedInUserID)

	commentID := c.Params("commentID")
	parsedCommentID, err := uuid.Parse(commentID)

	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var like models.UserProjectCommentLike
	if err := initializers.DB.Where("user_id=? AND project_comment_id=?", userID, parsedCommentID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserProjectCommentLike{
				ProjectCommentID: parsedCommentID,
				UserID:           userID,
			}
			result := initializers.DB.Create(&likeModel)
			if result.Error != nil {
				return &fiber.Error{Code: 500, Message: "Internal Server Error while adding the like."}
			}
			go routines.IncrementProjectCommentLikes(parsedCommentID, userID)
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}
		go routines.DecrementProjectCommentLikes(parsedCommentID)
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Liked/Unliked.",
	})
}
