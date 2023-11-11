package controllers

import (
	"errors"

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

func FollowUser(c *fiber.Ctx) error {
	loggedInUserIDStr := c.GetRespHeader("loggedInUserID")
	toFollowIDStr := c.Params("userID")

	loggedInUserID, _ := uuid.Parse(loggedInUserIDStr)
	toFollowID, err := uuid.Parse(toFollowIDStr)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	if loggedInUserID == toFollowID {
		return &fiber.Error{Code: 400, Message: "Cannot Follow Yourself."}
	}

	var follow models.FollowFollower
	if err := initializers.DB.Where("follower_id = ? AND followed_id = ?", loggedInUserID, toFollowID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var newFollow models.FollowFollower
			newFollow.FollowerID = loggedInUserID
			newFollow.FollowedID = toFollowID

			if err := initializers.DB.Create(&newFollow).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			go routines.IncrementCountsAndSendNotification(loggedInUserID, toFollowID) //! Check if race conditions of - + numFollowers happens when handling multiple users

			return c.Status(200).JSON(fiber.Map{
				"status":  "success",
				"message": "User followed successfully.",
			})
		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		return &fiber.Error{Code: 400, Message: "You are already following this user."}
	}
}

func UnfollowUser(c *fiber.Ctx) error {
	loggedInUserIDStr := c.GetRespHeader("loggedInUserID")
	toUnFollowIDStr := c.Params("userID")

	loggedInUserID, _ := uuid.Parse(loggedInUserIDStr)
	toUnFollowID, err := uuid.Parse(toUnFollowIDStr)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var follow models.FollowFollower
	if err := initializers.DB.Where("follower_id = ? AND followed_id = ?", loggedInUserID, toUnFollowID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "You do not follow this user."}
		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := initializers.DB.Where(&follow).Delete(&follow).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		go routines.DecrementCounts(loggedInUserID, toUnFollowID)

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User followed unfollowed.",
		})
	}
}

func RemoveFollow(c *fiber.Ctx) error {
	loggedInUserIDStr := c.GetRespHeader("loggedInUserID")
	followerToRemoveIDStr := c.Params("userID")

	loggedInUserID, _ := uuid.Parse(loggedInUserIDStr)
	followerToRemoveID, err := uuid.Parse(followerToRemoveIDStr)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var follow models.FollowFollower
	if err := initializers.DB.Where("follower_id = ? AND followed_id = ?", followerToRemoveID, loggedInUserID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "This user does not follow you."}
		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := initializers.DB.Where(&follow).Delete(&follow).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		go routines.DecrementCounts(followerToRemoveID, loggedInUserID)

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User followed removed from followers.",
		})
	}
}

func GetFollowers(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)
	searchDB := API.Search(c, 0)(paginatedDB)

	var followers []models.FollowFollower
	if err := searchDB.Preload("Follower").Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var followerUsers []models.User
	for _, follower := range followers {
		followerUsers = append(followerUsers, follower.Follower)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"users":   followerUsers,
	})
}

func GetFollowing(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)
	searchDB := API.Search(c, 0)(paginatedDB)

	var following []models.FollowFollower
	if err := searchDB.Preload("Followed").Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var followingUsers []models.User
	for _, followModel := range following {
		followingUsers = append(followingUsers, followModel.Followed)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"users":   followingUsers,
	})
}

func GetMyFollowing(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID, err := uuid.Parse(loggedInUserID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	var following []models.FollowFollower
	if err := initializers.DB.Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var followingUserIDs []uuid.UUID
	for _, followModel := range following {
		followingUserIDs = append(followingUserIDs, followModel.FollowedID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"userIDs": followingUserIDs,
	})
}

func GetMutuals(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	userID := c.Params("userID")

	joinQuery := initializers.DB.Joins("JOIN follow_followers AS ff1 ON ff1.follower_id = ? AND ff1.followed_id = follow_followers.follower_id", loggedInUserID).
		Joins("JOIN follow_followers AS ff2 ON ff2.follower_id = ? AND ff2.followed_id = follow_followers.followed_id", userID)

	var mutualFollowers []models.FollowFollower
	if err := joinQuery.
		Preload("Followed").
		Find(&mutualFollowers).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var count int64
	if err := joinQuery.
		Count(&count).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var mutuals []models.User
	for _, followObj := range mutualFollowers {
		mutuals = append(mutuals, followObj.Followed)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"users":   mutuals,
		"count":   count,
	})
}
