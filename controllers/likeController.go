package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
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

	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", parsedPostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

			post.NoLikes++

			if parsedLoggedInUserID != post.UserID {

				notification := models.Notification{
					NotificationType: 1,
					UserID:           post.UserID,
					SenderID:         parsedLoggedInUserID,
					PostID:           &post.ID,
				}

				if err := initializers.DB.Create(&notification).Error; err != nil {
					return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
				}
			}

		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}
		post.NoLikes--
	}

	result := initializers.DB.Save(&post)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the post."}
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

	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", parsedProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

			notification := models.Notification{
				NotificationType: 3,
				UserID:           project.UserID,
				SenderID:         userID,
				ProjectID:        &project.ID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
			}

			project.NoLikes++
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}

		project.NoLikes--
	}

	result := initializers.DB.Save(&project)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the project."}
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

	var comment models.PostComment
	if err := initializers.DB.First(&comment, "id = ?", parsedCommentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

			comment.NoLikes++
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		result := initializers.DB.Delete(&like)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}

		comment.NoLikes--
	}

	result := initializers.DB.Save(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the comment."}
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

	var comment models.ProjectComment
	if err := initializers.DB.First(&comment, "id = ?", parsedCommentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

			comment.NoLikes++
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}

	} else {
		result := initializers.DB.Delete(&like)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
		}

		comment.NoLikes--
	}

	result := initializers.DB.Save(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the comment."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Liked/Unliked.",
	})
}
