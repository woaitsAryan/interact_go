package controllers

import (
	"database/sql"

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

func GetComments(commentType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		itemID := c.Params(commentType + "ID")

		parsedItemID, err := uuid.Parse(itemID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID"}
		}

		paginatedDB := API.Paginator(c)(initializers.DB)
		db := paginatedDB.Preload("User").Where("is_flagged = ?", false)

		var comments []models.Comment
		switch commentType {
		case "post":
			db = db.Where("post_id = ? AND is_replied_comment = ?", parsedItemID, false)
		case "project":
			db = db.Where("project_id = ? AND is_replied_comment = ?", parsedItemID, false)
		case "event":
			db = db.Where("event_id = ? AND is_replied_comment = ?", parsedItemID, false)
		case "announcement":
			db = db.Where("announcement_id = ? AND is_replied_comment = ?", parsedItemID, false)
		case "comment":
			db = db.Where("parent_comment_id = ? AND is_replied_comment = ?", parsedItemID, true)
		}

		if err := db.Order("created_at DESC").Find(&comments).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":   "success",
			"comments": comments,
		})
	}
}

func AddComment(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content        string `json:"content"`
		PostID         string `json:"postID"`
		ProjectID      string `json:"projectID"`
		EventID        string `json:"eventID"`
		AnnouncementID string `json:"announcementID"`
		CommentID      string `json:"commentID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	postID := reqBody.PostID
	projectID := reqBody.ProjectID
	eventID := reqBody.EventID
	announcementID := reqBody.AnnouncementID
	commentID := reqBody.CommentID

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
	} else if announcementID != "" {
		parsedAnnouncementID, err := uuid.Parse(announcementID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}
		comment.AnnouncementID = &parsedAnnouncementID
		go routines.IncrementAnnouncementCommentsAndSendNotification(parsedAnnouncementID, parsedUserID)

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
	} else if commentID != "" {
		parsedCommentID, err := uuid.Parse(commentID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}
		comment.ParentCommentID = &parsedCommentID
		comment.IsRepliedComment = true

		var level int

		query := `
			WITH RECURSIVE comment_levels AS (
				SELECT id, parent_comment_id, 1 AS level
				FROM comments
				WHERE id = ?

				UNION ALL

				SELECT c.id, c.parent_comment_id, cl.level + 1
				FROM comments c
				INNER JOIN comment_levels cl ON c.id = cl.parent_comment_id
			)
			SELECT level
			FROM comment_levels
			ORDER BY level DESC
			LIMIT 1;
		`

		row := initializers.DB.Raw(query, parsedCommentID).Row()
		if err := row.Scan(&level); err != nil {
			if err == sql.ErrNoRows {
				return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		if level >= 5 {
			return &fiber.Error{Code: 400, Message: "Cannot reply to comments of level 5."}
		}

		comment.Level = level + 1
	}

	result := initializers.DB.Create(&comment)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	if commentID != "" {
		go routines.IncrementCommentReplies(uuid.MustParse(commentID))
	}

	if err := initializers.DB.Preload("User").First(&comment).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.CheckFlagComment(&comment)

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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.CheckFlagComment(&comment)

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
	if err := initializers.DB.
		Preload("Post").
		Preload("Project").
		Preload("Event").
		First(&comment, "id = ?", parsedCommentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Comment of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
	parentCommentID := comment.ParentCommentID

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if postID != nil {
		go routines.DecrementPostComments(*postID)
	} else if projectID != nil {
		go routines.DecrementProjectComments(*projectID)
	} else if eventID != nil {
		go routines.DecrementEventComments(*eventID)
	} else if parentCommentID != nil {
		go routines.DecrementCommentReplies(*parentCommentID)
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment deleted successfully",
	})
}
