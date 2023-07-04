package controllers

import (
	"errors"
	"reflect"
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	initializers.DB.Find(&users)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"users":   users,
	})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.Preload("Achievements").First(&user, "id = ?", userID)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	var user models.User
	initializers.DB.Preload("Achievements").First(&user, "id = ?", userID)
	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user of this ID found."}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")

	if user.ID.String() != loggedInUserID {
		utils.UpdateProfileViews(&user)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
	})
}

func UpdateMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var reqBody schemas.UserUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	picName, err := utils.SaveFile(c, "profilePic", "users/profilePics", true, 500, 500)
	if err != nil {
		return err
	}
	reqBody.ProfilePic = picName

	coverName, err := utils.SaveFile(c, "coverPic", "users/coverPics", true, 900, 400)
	if err != nil {
		return err
	}
	reqBody.CoverPic = coverName

	updateUserValue := reflect.ValueOf(&reqBody).Elem()
	userValue := reflect.ValueOf(&user).Elem()

	for i := 0; i < updateUserValue.NumField(); i++ {
		field := updateUserValue.Type().Field(i)
		fieldName := field.Name

		if fieldValue := updateUserValue.Field(i); fieldValue.IsValid() && fieldValue.String() != "" {
			userField := userValue.FieldByName(fieldName)
			if userField.IsValid() && userField.CanSet() {
				userField.Set(fieldValue)
			}
		}
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func DeleteMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	user.Active = false

	if err := initializers.DB.Save(&user).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

func UpdatePassword(c *fiber.Ctx) error {

	var reqBody struct {
		Password        string `json:"password"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	if reqBody.NewPassword != reqBody.ConfirmPassword {
		return &fiber.Error{Code: 400, Message: "Passwords do not match."}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.First(&user, "id = ?", loggedInUserID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorret Password."}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), 10)

	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}

	user.Password = string(hash)
	user.PasswordChangedAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Password updated successfully",
	})
}
