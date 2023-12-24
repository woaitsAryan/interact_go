package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
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

	var comments []models.Comment
	if err := paginatedDB.Preload("User").Where("post_id=?", parsedPostID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
	if err := paginatedDB.Preload("User").Where("project_id=?", parsedProjectID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"comments": comments,
	})
}

func GetEventComments(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	parsedEventID, err := uuid.Parse(eventID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var comments []models.Comment
	if err := paginatedDB.Preload("User").Where("event_id=?", parsedEventID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		EventID   string `json:"eventID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	postID := reqBody.PostID
	projectID := reqBody.ProjectID
	eventID := reqBody.EventID

	comment := models.Comment{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	if postID != "" {
		parsedPostID, err := uuid.Parse(postID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}
		comment.PostID = &parsedPostID
		go routines.IncrementPostCommentsAndSendNotification(parsedPostID, parsedUserID)
	} else if projectID != "" {
		parsedProjectID, err := uuid.Parse(projectID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}
		comment.ProjectID = &parsedProjectID
		go routines.IncrementProjectCommentsAndSendNotification(parsedProjectID, parsedUserID)
	} else if eventID != "" {
		parsedEventID, err := uuid.Parse(eventID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}
		comment.EventID = &parsedEventID
		go routines.IncrementEventCommentsAndSendNotification(parsedEventID, parsedUserID)
	}

	result := initializers.DB.Create(&comment)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	if err := initializers.DB.Preload("User").First(&comment).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment Added",
		"comment": comment,
	})
}

func UpdateComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	commentID := c.Params("commentID")

	parsedCommentID, err := uuid.Parse(commentID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var comment models.Comment
	if err := initializers.DB.First(&comment, "id = ? AND user_id=?", parsedCommentID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

func DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("commentID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedCommentID, err := uuid.Parse(commentID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var comment models.Comment
	if err := initializers.DB.Preload("Post").
		Preload("Project").
		Preload("Event").
		First(&comment, "id = ?", parsedCommentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if comment.UserID != parsedLoggedInUserID {
		if comment.PostID != nil {
			if comment.Post.UserID != parsedLoggedInUserID {
				return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
			}
		} else if comment.ProjectID != nil {
			if comment.Project.UserID != parsedLoggedInUserID {
				return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
			}
		} else if comment.EventID != nil {
			if comment.Event.OrganizationID != parsedLoggedInUserID {
				return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
			}
		} else {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
	}

	postID := comment.PostID
	projectID := comment.ProjectID
	eventID := comment.EventID

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if postID != nil {
		go routines.DecrementPostComments(*postID)
	} else if projectID != nil {
		go routines.DecrementProjectComments(*projectID)
	} else if eventID != nil {
		go routines.DecrementEventComments(*eventID)
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment deleted successfully",
	})
}
