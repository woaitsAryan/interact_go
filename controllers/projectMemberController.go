package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddMember(c *fiber.Ctx) error {

	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID"}
	}

	var reqBody struct {
		UserID string
		Title  string
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", reqBody.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No User of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	membership := models.Membership{
		ProjectID: parsedProjectID,
		UserID:    user.ID,
		Role:      "",
		Title:     reqBody.Title,
	}

	result := initializers.DB.Create(&membership)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating membership."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User added to the project.",
	})
}

func RemoveMember(c *fiber.Ctx) error {

	membershipID := c.Params("membershipID")

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var membership models.Membership
	if err := initializers.DB.First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&membership)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting membership."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User removed to the project.",
	})
}

func ChangeMemberRole(c *fiber.Ctx) error {

	membershipID := c.Params("membershipID")

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var reqBody struct {
		Role string
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var membership models.Membership
	if err := initializers.DB.First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	membership.Role = reqBody.Role

	result := initializers.DB.Save(&membership)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating membership."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User membership updated.",
	})
}
