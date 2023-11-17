package controllers

import (
	"errors"
	"strconv"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func EditProfile(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	var profile models.Profile
	if err := initializers.DB.First(&profile, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			newProfile := models.Profile{
				UserID: parsedUserID,
			}

			result := initializers.DB.Create(&newProfile)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			return &fiber.Error{Code: 400, Message: "Some Error Occurred, Please Try Again."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.ProfileUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if reqBody.School != nil {
		profile.School = *reqBody.School
	}
	if reqBody.Description != nil {
		profile.Description = *reqBody.Description
	}
	if reqBody.Areas != nil {
		profile.AreasOfCollaboration = *reqBody.Areas
	}
	if reqBody.Degree != nil {
		profile.Degree = *reqBody.Degree
	}
	if reqBody.Hobbies != nil {
		profile.Hobbies = *reqBody.Hobbies
	}
	if reqBody.YOG != nil {
		year, err := strconv.Atoi(*reqBody.YOG)
		if err == nil {
			profile.YearOfGraduation = year
		}
	}

	result := initializers.DB.Save(&profile)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Profile Edited.",
		"profile": profile,
	})
}

func AddAchievement(c *fiber.Ctx) error {
	// userID := c.GetRespHeader("loggedInUserID")
	// parsedUserID, _ := uuid.Parse(userID)

	var reqBody schemas.AchievementCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	for _, achievement := range reqBody.Achievements {

		var achievementModel models.Achievement
		// achievementModel.UserID = parsedUserID

		if achievement.ID == "" {
			achievementModel.Title = achievement.Title
			achievementModel.Skills = achievement.Skills
			err := initializers.DB.Create(&achievementModel).Error

			if err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
		} else {
			err := initializers.DB.First(&achievementModel, "id = ?", achievement.ID).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &fiber.Error{Code: 400, Message: "Invalid ID."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			} else {
				achievementModel.Skills = achievement.Skills
				achievementModel.Title = achievement.Title
				if err := initializers.DB.Save(&achievementModel).Error; err != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if err := initializers.DB.Delete(&achievement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement deleted successfully",
	})
}
