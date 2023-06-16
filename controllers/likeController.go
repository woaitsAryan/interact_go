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
	userID, _ := uuid.Parse(loggedInUserID)

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
	if err := initializers.DB.Where("user_id=? AND post_id=?", userID, parsedPostID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserPostLike{
				PostID: parsedPostID,
				UserID: userID,
			}

			result := initializers.DB.Create(&likeModel)

			if result.Error != nil {
				return &fiber.Error{Code: 500, Message: "Internal Server Error while adding the like."}
			}

			post.NoLikes++
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&like)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
	}

	post.NoLikes--

	result = initializers.DB.Save(&post)

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

			project.NoLikes++
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&like)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
	}

	project.NoLikes--

	result = initializers.DB.Save(&project)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the project."}
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
	if err := initializers.DB.First(&comment, "id = ?", parsedCommentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var like models.UserProjectLike
	if err := initializers.DB.Where("user_id=? AND comment_id=?", userID, parsedCommentID).First(&like).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.UserCommentLike{
				CommentID: parsedCommentID,
				UserID:    userID,
			}

			result := initializers.DB.Create(&likeModel)

			if result.Error != nil {
				return &fiber.Error{Code: 500, Message: "Internal Server Error while adding the like."}
			}

			comment.NoLikes++
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&like)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the like."}
	}

	comment.NoLikes--

	result = initializers.DB.Save(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the comment."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Liked/Unliked.",
	})
}
