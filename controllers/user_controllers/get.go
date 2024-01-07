package user_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.Preload("Profile").First(&user, "username = ?", username)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user of this username found."}
	}

	if user.ID.String() != loggedInUserID {
		routines.UpdateProfileViews(&user)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
	})
}

func GetMyLikes(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var likes []models.Like
	if err := initializers.DB.
		Find(&likes, "user_id = ?", loggedInUserID).Error; err != nil {
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
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"likes":   likeIDs,
	})
}

func GetMyOrgMemberships(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	populate := c.Query("populate", "false")

	var memberships []models.OrganizationMembership

	if populate == "true" {
		if err := initializers.DB.
			Preload("Organization").
			Preload("Organization.User").
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
