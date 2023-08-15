package organization_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddMember(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	organizationID := c.Params("organizationID")

	parsedOrganizationID, err := uuid.Parse(organizationID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	var reqBody struct {
		UserID string
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", reqBody.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No User of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "id = ? and user_id=?", parsedOrganizationID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.UserID == organization.UserID.String() {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.Where("user_id=? AND organization_id=?", user.ID, parsedOrganizationID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// var existingInvitation models.ProjectInvitation
			// err := initializers.DB.Where("user_id=? AND project_id=? AND status=0", user.ID, parsedProjectID).First(&existingInvitation).Error
			// if err == nil {
			// 	return &fiber.Error{Code: 400, Message: "Have already invited this User."}
			// }

			// var invitation models.ProjectInvitation
			// invitation.ProjectID = parsedProjectID
			// invitation.UserID = user.ID
			// invitation.Title = reqBody.Title

			// result := initializers.DB.Create(&invitation)

			// if result.Error != nil {
			// 	return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			// }

			// invitation.User = user

			return c.Status(201).JSON(fiber.Map{
				"status":  "success",
				"message": "Invitation sent to the user.",
			})
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	} else {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}
}

func RemoveMember(c *fiber.Ctx) error {
	membershipID := c.Params("membershipID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.Preload("Organization").First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if membership.Organization.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
	}

	if membership.Role == "Owner" {
		return &fiber.Error{Code: 403, Message: "Invalid Route to perform this action."}
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

func LeaveOrganization(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var membership models.OrganizationMembership
	if err := initializers.DB.Preload("Organization").First(&membership, "user_id=? AND organization_id = ?", loggedInUserID, projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	result := initializers.DB.Delete(&membership)
	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting membership."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "You left the organization",
	})
}

func ChangeMemberRole(c *fiber.Ctx) error {
	membershipID := c.Params("membershipID")

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var reqBody struct {
		Role models.OrganizationRole
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.Role != models.Member && reqBody.Role != models.Manager {
		return &fiber.Error{Code: 403, Message: "Invalid route to perform this action."}
	}

	membership.Role = reqBody.Role

	result := initializers.DB.Save(&membership)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User membership updated.",
	})
}
