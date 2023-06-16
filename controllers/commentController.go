package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPostComments(c *fiber.Ctx) error {
	postID := c.Params("postID")

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var comments []models.Comment
	if err := paginatedDB.Where("post_id=?", parsedPostID).Find(&comments).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"comments": comments,
	})
}

func GetProjectComments(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var comments []models.Comment
	if err := paginatedDB.Where("project_id=?", parsedProjectID).Find(&comments).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"comments": comments,
	})
}

func AddComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content   string `json:"content"`
		PostID    string `json:"postID"`
		ProjectID string `json:"projectID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	postID := reqBody.PostID
	projectID := reqBody.ProjectID

	comment := models.Comment{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	notification := models.Notification{
		SenderID: parsedUserID,
	}

	if postID != "" {
		parsedPostID, err := uuid.Parse(postID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		var post models.Post
		if err := initializers.DB.First(&post, "id=?", parsedPostID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}

		comment.PostID = parsedPostID
		notification.NotificationType = 2
		notification.UserID = post.UserID
		notification.PostID = post.ID

	} else if projectID != "" {
		parsedProjectID, err := uuid.Parse(projectID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		var project models.Project
		if err := initializers.DB.First(&project, "id=?", parsedProjectID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}

		comment.ProjectID = parsedProjectID
		notification.NotificationType = 4
		notification.UserID = project.UserID
		notification.PostID = project.ID

	} else {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	if err := initializers.DB.Create(&notification).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
	}

	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the comment."}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Added",
	})
}

func UpdateComment(c *fiber.Ctx) error {
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

	var reqBody struct {
		Content string `json:"comment"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if reqBody.Content != "" {
		comment.Content = reqBody.Content
	}

	comment.Edited = true

	if err := initializers.DB.Save(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

func DeleteComemnt(c *fiber.Ctx) error {
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

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment deleted successfully",
	})
}
