package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddMember(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
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

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? and user_id=?", parsedProjectID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if reqBody.UserID == project.UserID.String() {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}

	var membership models.Membership
	if err := initializers.DB.Where("user_id=? AND project_id=?", user.ID, parsedProjectID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			var existingInvitation models.ProjectInvitation
			err := initializers.DB.Where("user_id=? AND project_id=? AND status=0", user.ID, parsedProjectID).First(&existingInvitation).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "Have already invited this User."}
			}

			var invitation models.ProjectInvitation
			invitation.ProjectID = parsedProjectID
			invitation.UserID = user.ID
			invitation.Title = reqBody.Title

			result := initializers.DB.Create(&invitation)

			if result.Error != nil {
				return &fiber.Error{Code: 500, Message: "Internal Server Error while creating Invitation."}
			}

			invitation.User = user

			return c.Status(201).JSON(fiber.Map{
				"status":     "success",
				"message":    "Invitation sent to the user.",
				"invitation": invitation,
			})
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	} else {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}
}

func RemoveMember(c *fiber.Ctx) error { //! Only project creator can access

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

func ChangeMemberRole(c *fiber.Ctx) error { //! Only project creator can access

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
