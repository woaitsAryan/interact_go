package controllers

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddAchievement(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var reqBody schemas.AchievementCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	for _, achievement := range reqBody.Achievements {

		var achievementModel models.Achievement
		achievementModel.UserID = parsedUserID

		if achievement.ID == "" {
			achievementModel.Title = achievement.Title
			achievementModel.Skills = achievement.Skills
			err := initializers.DB.Create(&achievementModel).Error

			if err != nil {
				return &fiber.Error{Code: 500, Message: "Database Error While creating achievement."}
			}
		} else {
			err := initializers.DB.First(&achievementModel, "id = ?", achievement.ID).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &fiber.Error{Code: 400, Message: "Invalid ID."}
				}
				return &fiber.Error{Code: 500, Message: "Database Error."}
			} else {
				achievementModel.Skills = achievement.Skills
				achievementModel.Title = achievement.Title
				if err := initializers.DB.Save(&achievementModel).Error; err != nil {
					return &fiber.Error{Code: 500, Message: "Database Error."}
				}
			}
		}

	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement added successfully",
	})
}

func DeleteAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("achievementID")
	userID := c.GetRespHeader("loggedInUserID")

	var achievement models.Achievement
	if err := initializers.DB.Where("user_id=? AND id=?", userID, achievementID).First(&achievement).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Achievement of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&achievement).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement deleted successfully",
	})
}
