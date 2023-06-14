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
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	toFollowID := c.Params("userID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", loggedInUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "Server error. Log in again."}
		}
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	var toFollow models.User
	if err := initializers.DB.First(&toFollow, "id = ?", toFollowID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	var existingFollower models.User
	initializers.DB.Model(&toFollow).Association("Followers").Find(&existingFollower, "id = ?", user.ID)

	if existingFollower.ID != uuid.Nil {
		return &fiber.Error{Code: 400, Message: "You are already following this user."}
	}

	user.Following = append(user.Following, &toFollow)
	if err := initializers.DB.Save(&user).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	toFollow.Followers = append(toFollow.Followers, &user)
	if err := initializers.DB.Save(&toFollow).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User followed successfully.",
	})
}

func UnfollowUser(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	toUnfollowID := c.Params("userID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", loggedInUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "Server error. Log in again."}
		}
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	var toUnfollow models.User
	if err := initializers.DB.First(&toUnfollow, "id = ?", toUnfollowID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	var existingFollower models.User
	initializers.DB.Model(&toUnfollow).Association("Followers").Find(&existingFollower, "id = ?", user.ID)

	if existingFollower.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "You are not following this user."}
	}

	initializers.DB.Model(&user).Association("Following").Delete(&toUnfollow)

	initializers.DB.Model(&toUnfollow).Association("Followers").Delete(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User unfollowed successfully.",
	})
}

func RemoveFollow(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	followerToRemoveID := c.Params("userID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", loggedInUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "Server error. Log in again."}
		}
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	var follower models.User
	if err := initializers.DB.First(&follower, "id = ?", followerToRemoveID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No follower of this ID found."}
		}
		return &fiber.Error{Code: 400, Message: "Database Error."}
	}

	var existingFollowing models.User
	initializers.DB.Model(&user).Association("Following").Find(&existingFollowing, "id = ?", follower.ID)

	if existingFollowing.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "This follower is not in your following list."}
	}

	initializers.DB.Model(&user).Association("Following").Delete(&follower)

	initializers.DB.Model(&follower).Association("Followers").Delete(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Follower removed successfully.",
	})
}
