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

func LikeItem(likeType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

		var like models.Like

		switch likeType {
		case "post":
			postID := c.Params("postID")
			parsedPostID, err := uuid.Parse(postID)

			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID"}
			}

			err = initializers.DB.Where("user_id=? AND post_id=?", parsedLoggedInUserID, parsedPostID).First(&like).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					likeModel := models.Like{
						PostID: &parsedPostID,
						UserID: parsedLoggedInUserID,
					}

					result := initializers.DB.Create(&likeModel)
					if result.Error != nil {
						return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
					}
					go routines.IncrementPostLikesAndSendNotification(parsedPostID, parsedLoggedInUserID)

				} else {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}
			} else {
				result := initializers.DB.Delete(&like)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
				go routines.DecrementPostLikes(parsedPostID)
			}

		case "project":
			projectID := c.Params("projectID")
			parsedProjectID, err := uuid.Parse(projectID)

			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID"}
			}

			var like models.Like
			if err := initializers.DB.Where("user_id=? AND project_id=?", parsedLoggedInUserID, parsedProjectID).First(&like).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					likeModel := models.Like{
						ProjectID: &parsedProjectID,
						UserID:    parsedLoggedInUserID,
					}

					result := initializers.DB.Create(&likeModel)

					if result.Error != nil {
						return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
					}
					go routines.IncrementProjectLikesAndSendNotification(parsedProjectID, parsedLoggedInUserID)
				} else {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}
			} else {
				result := initializers.DB.Delete(&like)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
				go routines.DecrementProjectLikes(parsedProjectID)
			}
		case "comment":
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
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			var like models.Like
			if err := initializers.DB.Where("user_id=? AND comment_id=?", parsedLoggedInUserID, parsedCommentID).First(&like).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					likeModel := models.Like{
						CommentID: &comment.ID,
						UserID:    parsedLoggedInUserID,
					}
					result := initializers.DB.Create(&likeModel)
					if result.Error != nil {
						return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
					}
					go routines.IncrementCommentLikes(parsedCommentID, parsedLoggedInUserID)
				} else {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}
			} else {
				result := initializers.DB.Delete(&like)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
				go routines.DecrementCommentLikes(parsedCommentID)
			}
		case "event":
			eventID := c.Params("eventID")
			parsedEventID, err := uuid.Parse(eventID)

			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID"}
			}

			var like models.Like
			if err := initializers.DB.Where("user_id=? AND event_id=?", parsedLoggedInUserID, parsedEventID).First(&like).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					likeModel := models.Like{
						EventID: &parsedEventID,
						UserID:  parsedLoggedInUserID,
					}
					result := initializers.DB.Create(&likeModel)
					if result.Error != nil {
						return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
					}
					go routines.IncrementEventLikesAndSendNotification(parsedEventID, parsedLoggedInUserID)
				} else {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}
			} else {
				result := initializers.DB.Delete(&like)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
				go routines.DecrementEventLikes(parsedEventID)
			}
		case "review":
			reviewID := c.Params("reviewID")
			parsedReviewID, err := uuid.Parse(reviewID)

			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID"}
			}

			var like models.Like
			if err := initializers.DB.Where("user_id=? AND review_id=?", parsedLoggedInUserID, parsedReviewID).First(&like).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					likeModel := models.Like{
						EventID: &parsedReviewID,
						UserID:  parsedLoggedInUserID,
					}
					result := initializers.DB.Create(&likeModel)
					if result.Error != nil {
						return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
					}
					go routines.IncrementReviewLikes(parsedReviewID, parsedLoggedInUserID)
				} else {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}
			} else {
				result := initializers.DB.Delete(&like)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
				go routines.DecrementReviewLikes(parsedReviewID)
			}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Liked/Unliked.",
		})
	}
}
