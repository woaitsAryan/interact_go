package user_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetViews(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	viewsArr, count, err := utils.GetProfileViews(parsedUserID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"viewsArr": viewsArr,
		"count":    count,
	})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.
		Preload("Profile").
		Preload("Profile.Achievements").
		First(&user, "id = ?", userID)

	var profile models.Profile
	if err := initializers.DB.First(&profile, "user_id = ?", userID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
		"profile": profile,
	})
}

func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	//TODO20 add error handing here
	var user models.User
	initializers.DB.Preload("Profile").First(&user, "username = ?", username)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user of this username found."}
	}

	var profile models.Profile
	if err := initializers.DB.First(&profile, "user_id = ?", user.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if user.ID.String() != loggedInUserID {
		routines.UpdateProfileViews(&user)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
		"profile": profile,
	})
}

func GetMyLikes(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var likes []models.Like
	if err := initializers.DB.
		Find(&likes, "user_id = ? AND status = 0", loggedInUserID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var likeIDs []string
	for _, like := range likes {
		if like.PostID != nil {
			likeIDs = append(likeIDs, like.PostID.String())
		} else if like.ProjectID != nil {
			likeIDs = append(likeIDs, like.ProjectID.String())
		} else if like.CommentID != nil {
			likeIDs = append(likeIDs, like.CommentID.String())
		} else if like.EventID != nil {
			likeIDs = append(likeIDs, like.EventID.String())
		} else if like.ReviewID != nil {
			likeIDs = append(likeIDs, like.ReviewID.String())
		} else if like.AnnouncementID != nil {
			likeIDs = append(likeIDs, like.AnnouncementID.String())
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"likes":   likeIDs,
	})
}

func GetMyDislikes(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var likes []models.Like
	if err := initializers.DB.
		Find(&likes, "user_id = ? AND status = -1", loggedInUserID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var likeIDs []string
	for _, like := range likes {
		if like.PostID != nil {
			likeIDs = append(likeIDs, like.PostID.String())
		} else if like.ProjectID != nil {
			likeIDs = append(likeIDs, like.ProjectID.String())
		} else if like.CommentID != nil {
			likeIDs = append(likeIDs, like.CommentID.String())
		} else if like.EventID != nil {
			likeIDs = append(likeIDs, like.EventID.String())
		} else if like.ReviewID != nil {
			likeIDs = append(likeIDs, like.ReviewID.String())
		} else if like.AnnouncementID != nil {
			likeIDs = append(likeIDs, like.AnnouncementID.String())
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "User Found",
		"dislikes": likeIDs,
	})
}

func GetMyOrgMemberships(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	populate := c.Query("populate", "false")

	var memberships []models.OrganizationMembership

	if populate == "true" {
		if err := initializers.DB.
			Preload("Organization").
			Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
				return db.Select(select_fields.User)
			}).
			Find(&memberships, "user_id = ?", loggedInUserID).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	} else {
		if err := initializers.DB.
			Find(&memberships, "user_id = ?", loggedInUserID).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "User Found",
		"memberships": memberships,
	})
}

func GetMyVotedOptions(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var options []models.Option
	if err := initializers.DB.
		Joins("JOIN voted_by ON voted_by.option_id = options.id").
		Where("voted_by.user_id = ?", loggedInUserID).
		Find(&options).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var optionIDs []string
	for _, option := range options {
		optionIDs = append(optionIDs, option.ID.String())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"options": optionIDs,
	})
}
