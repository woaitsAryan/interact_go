package controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("loggedInUser")
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	var user models.User
	initializers.DB.First(&user, "id = ?", userID)
	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user of this ID found."}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")

	if user.ID.String() != loggedInUserID {
		// Updating Profile Views
		today := time.Now().UTC().Truncate(24 * time.Hour)
		var profileView models.ProfileView
		initializers.DB.Where("user_id = ? AND date = ?", user.ID, today).First(&profileView)

		if profileView.ID == uuid.Nil {
			profileView = models.ProfileView{
				UserID: user.ID,
				Date:   today,
				Count:  1,
			}
			initializers.DB.Create(&profileView)
		} else {
			profileView.Count++
			initializers.DB.Save(&profileView)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
	})
}
