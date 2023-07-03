package controllers

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
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

	var toFollowUser models.User
	err := initializers.DB.First(&toFollowUser, "id=?", toFollowID).Error
	if err != nil {
		return &fiber.Error{Code: 500, Message: "No User with this ID exists."}
	}

	var loggedInUser models.User
	initializers.DB.First(&loggedInUser, "id=?", loggedInUserID)

	var follow models.FollowFollower
	if err := initializers.DB.Where("follower_id = ? AND followed_id = ?", loggedInUserID, toFollowID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var newFollow models.FollowFollower
			newFollow.FollowerID = loggedInUserID
			newFollow.FollowedID = toFollowID

			if err := initializers.DB.Create(&newFollow).Error; err != nil {
				return &fiber.Error{Code: 500, Message: "Database Error while creating follow."}
			}

			notification := models.Notification{
				NotificationType: 0,
				UserID:           toFollowUser.ID,
				SenderID:         loggedInUserID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				return &fiber.Error{Code: 500, Message: "Database Error while creating notification."}
			}

			toFollowUser.NoFollowers++
			if err := initializers.DB.Save(&toFollowUser).Error; err != nil {
				return &fiber.Error{Code: 500, Message: "Database Error while incrementing number followers."}
			}

			loggedInUser.NoFollowing++
			if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
				return &fiber.Error{Code: 500, Message: "Database Error while incrementing number following."}
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

		var toUnFollowUser models.User
		initializers.DB.First(&toUnFollowUser, "id=?", toUnFollowID)

		var loggedInUser models.User
		initializers.DB.First(&loggedInUser, "id=?", loggedInUserID)

		toUnFollowUser.NoFollowers--
		if err := initializers.DB.Save(&toUnFollowUser).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while decrementing number followers."}
		}

		loggedInUser.NoFollowing--
		if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while decrementing number following."}
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

		var followerToRemove models.User
		initializers.DB.First(&followerToRemove, "id=?", followerToRemoveID)

		var loggedInUser models.User
		initializers.DB.First(&loggedInUser, "id=?", loggedInUserID)

		followerToRemove.NoFollowing--
		if err := initializers.DB.Save(&followerToRemove).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while decrementing number following."}
		}

		loggedInUser.NoFollowers--
		if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while decrementing number followers."}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User followed removed from followers.",
		})
	}
}

func GetFollowers(c *fiber.Ctx) error { //! Add search here
	userIDStr := c.Params("userID")
	userID := uuid.MustParse(userIDStr)

	paginatedDB := API.Paginator(c)(initializers.DB)

	var followers []models.FollowFollower
	if err := paginatedDB.Preload("Follower").Select("follower_id").Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
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

func GetFollowing(c *fiber.Ctx) error { //! Add search here
	userIDStr := c.Params("userID")
	userID := uuid.MustParse(userIDStr)

	paginatedDB := API.Paginator(c)(initializers.DB)

	var following []models.FollowFollower
	if err := paginatedDB.Preload("Followed").Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"following": following,
	})
}
