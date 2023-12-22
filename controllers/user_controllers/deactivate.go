package user_controllers

import (
	"errors"
	"time"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Deactivate(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		VerificationCode string `json:"otp"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	data, err := cache.GetOtpFromCache(user.ID.String())
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data), []byte(reqBody.VerificationCode)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorrect OTP"}
	}

	user.Active = false
	user.DeactivatedAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Account Deactivated",
	})
}
