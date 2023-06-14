package controllers

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FollowUser(c *fiber.Ctx) error {
	loggedInUserIDStr := c.GetRespHeader("loggedInUserID")
	toFollowIDStr := c.Params("userID")

	loggedInUserID := uuid.MustParse(loggedInUserIDStr)
	toFollowID := uuid.MustParse(toFollowIDStr)

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
				return &fiber.Error{Code: 500, Message: "Database Error."}
			}

			return c.Status(200).JSON(fiber.Map{
				"status":  "success",
				"message": "User followed successfully.",
			})
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		return &fiber.Error{Code: 400, Message: "You are already following this user."}
	}

}

func UnfollowUser(c *fiber.Ctx) error {
	loggedInUserIDStr := c.GetRespHeader("loggedInUserID")
	toUnFollowIDStr := c.Params("userID")

	loggedInUserID := uuid.MustParse(loggedInUserIDStr)
	toUnFollowID := uuid.MustParse(toUnFollowIDStr)

	var follow models.FollowFollower
	if err := initializers.DB.Where("follower_id = ? AND followed_id = ?", loggedInUserID, toUnFollowID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "You do not follow this user."}
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		if err := initializers.DB.Where(&follow).Delete(&follow).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User followed unfollowed.",
		})
	}
}

func RemoveFollow(c *fiber.Ctx) error {
	loggedInUserIDStr := c.GetRespHeader("loggedInUserID")
	followerToRemoveIDStr := c.Params("userID")

	loggedInUserID := uuid.MustParse(loggedInUserIDStr)
	followerToRemoveID := uuid.MustParse(followerToRemoveIDStr)

	var follow models.FollowFollower
	if err := initializers.DB.Where("follower_id = ? AND followed_id = ?", followerToRemoveID, loggedInUserID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "This user does not follow you."}
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		if err := initializers.DB.Where(&follow).Delete(&follow).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User followed removed from followers.",
		})
	}
}

func GetFollowers(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	userID := uuid.MustParse(userIDStr)

	var followers []models.FollowFollower
	if err := initializers.DB.Preload("Follower").Select("follower_id").Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var followerUsers []models.User
	for _, follower := range followers {
		followerUsers = append(followerUsers, follower.Follower)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"followers": followerUsers,
	})
}

func GetFollowing(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	userID := uuid.MustParse(userIDStr)

	var following []models.FollowFollower
	if err := initializers.DB.Preload("Followed").Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"following": following,
	})
}
