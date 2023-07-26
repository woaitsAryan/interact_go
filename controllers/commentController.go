package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
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

	var comments []models.PostComment
	if err := paginatedDB.Preload("User").Where("post_id=?", parsedPostID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
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

	var comments []models.ProjectComment
	if err := paginatedDB.Preload("User").Where("project_id=?", parsedProjectID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"comments": comments,
	})
}

func AddPostComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content string `json:"content"`
		PostID  string `json:"postID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	postID := reqBody.PostID
	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	comment := models.PostComment{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	comment.PostID = parsedPostID
	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	if err := initializers.DB.Preload("User").First(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	go routines.IncrementPostCommentsAndSendNotification(parsedPostID, parsedUserID)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Added",
		"comment": comment,
	})
}

func AddProjectComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content   string `json:"content"`
		ProjectID string `json:"projectID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	projectID := reqBody.ProjectID

	comment := models.ProjectComment{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	comment.ProjectID = parsedProjectID

	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	if err := initializers.DB.Preload("User").First(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	go routines.IncrementProjectCommentsAndSendNotification(parsedProjectID, parsedUserID)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Added",
		"comment": comment,
	})
}

func UpdatePostComment(c *fiber.Ctx) error {
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
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
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
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

func DeletePostComment(c *fiber.Ctx) error {
	commentID := c.Params("commentID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedCommentID, err := uuid.Parse(commentID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var comment models.PostComment
	if err := initializers.DB.First(&comment, "id = ? AND user_id = ?", parsedCommentID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	go routines.DecrementPostComments(comment.PostID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment deleted successfully",
	})
}

func DeleteProjectComment(c *fiber.Ctx) error {
	commentID := c.Params("commentID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedCommentID, err := uuid.Parse(commentID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var comment models.ProjectComment
	if err := initializers.DB.First(&comment, "id = ? AND user_id = ?", parsedCommentID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	go routines.DecrementProjectComments(comment.ProjectID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment deleted successfully",
	})
}

func GetMyLikedPostsComments(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var postCommentLikes []models.UserPostCommentLike
	if err := initializers.DB.Where("user_id = ?", loggedInUserID).Find(&postCommentLikes).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	var postCommentIDs []string
	for _, postCommentLike := range postCommentLikes {
		postCommentIDs = append(postCommentIDs, postCommentLike.PostCommentID.String())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"comments": postCommentIDs,
	})
}

func GetMyLikedProjectsComments(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var projectCommentLikes []models.UserProjectCommentLike
	if err := initializers.DB.Where("user_id = ?", loggedInUserID).Find(&projectCommentLikes).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	var projectCommentIDs []string
	for _, projectCommentLike := range projectCommentLikes {
		projectCommentIDs = append(projectCommentIDs, projectCommentLike.ProjectCommentID.String())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"comments": projectCommentIDs,
	})
}
