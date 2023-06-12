package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user")
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
	if user.ID == 0 {
		return &fiber.Error{Code: 400, Message: "No user of this ID found."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}
