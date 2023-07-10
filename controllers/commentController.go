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

	var comments []models.PostComment
	if err := paginatedDB.Preload("User").Where("post_id=?", parsedPostID).Order("created_at DESC").Find(&comments).Error; err != nil {
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

	var comments []models.ProjectComment
	if err := paginatedDB.Preload("User").Where("project_id=?", parsedProjectID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
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

	var post models.Post
	if err := initializers.DB.First(&post, "id=?", parsedPostID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
	}

	comment.PostID = parsedPostID
	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the comment."}
	}

	post.NoComments++
	result = initializers.DB.Save(&post)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the post."}
	}

	if parsedUserID != post.UserID {
		notification := models.Notification{
			SenderID:         parsedUserID,
			NotificationType: 2,
			UserID:           post.UserID,
			PostID:           &post.ID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
		}
	}

	if err := initializers.DB.Preload("User").First(&comment).Error; err != nil {
		// Handle the error if the preload fails
		return &fiber.Error{Code: 500, Message: "Internal Server Error while loading the user."}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Added",
		"comment": comment, // The comment is updated with ids n everything
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

	notification := models.Notification{
		SenderID: parsedUserID,
	}

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
	notification.PostID = &project.ID

	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the comment."}
	}

	project.NoComments++
	result = initializers.DB.Save(&project)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the Project."}
	}

	// if err := initializers.DB.Create(&notification).Error; err != nil {
	// 	return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
	// }

	if err := initializers.DB.Preload("User").First(&comment).Error; err != nil {
		// Handle the error if the preload fails
		return &fiber.Error{Code: 500, Message: "Internal Server Error while loading the user."}
	}

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
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var post models.Post
	if err := initializers.DB.First(&post, "id=?", comment.PostID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	post.NoComments--
	result := initializers.DB.Save(&post)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the post."}
	}

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
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id=?", comment.ProjectID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	project.NoComments--
	result := initializers.DB.Save(&project)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while saving the project."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment deleted successfully",
	})
}

func GetMyLikedPostsComments(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var postCommentLikes []models.UserPostCommentLike
	if err := initializers.DB.Where("user_id = ?", loggedInUserID).Find(&postCommentLikes).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
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
		return &fiber.Error{Code: 500, Message: "Database Error."}
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
