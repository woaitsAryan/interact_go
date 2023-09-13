package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetNonMembers(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	var project models.Project
	if err := initializers.DB.Where("id = ?", projectID).Preload("Memberships").First(&project).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID"}
	}

	var membershipUserIDs []string

	for _, membership := range project.Memberships {
		membershipUserIDs = append(membershipUserIDs, membership.UserID.String())
	}

	membershipUserIDs = append(membershipUserIDs, project.UserID.String())

	searchedDB := API.Search(c, 0)(initializers.DB)

	var users []models.User
	if err := searchedDB.Where("id NOT IN (?)", membershipUserIDs).Limit(10).Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? and user_id=?", parsedProjectID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.UserID == project.UserID.String() {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}

	var membership models.Membership
	if err := initializers.DB.Where("user_id=? AND project_id=?", user.ID, parsedProjectID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			var existingInvitation models.Invitation
			err := initializers.DB.Where("user_id=? AND project_id=? AND status=0", user.ID, parsedProjectID).First(&existingInvitation).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "Have already invited this User."}
			}

			var existingApplication models.Application
			err = initializers.DB.Where("user_id=? AND project_id=? AND (status=0 OR status=1)", user.ID, parsedProjectID).First(&existingApplication).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "User has already applied for this project."}
			}

			var invitation models.Invitation
			invitation.ProjectID = &parsedProjectID
			invitation.UserID = user.ID
			invitation.Title = reqBody.Title

			result := initializers.DB.Create(&invitation)

			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			invitation.User = user

			return c.Status(201).JSON(fiber.Map{
				"status":     "success",
				"message":    "Invitation sent to the user.",
				"invitation": invitation,
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

	var membership models.Membership
	if err := initializers.DB.Preload("Project").First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if membership.Project.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
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

func LeaveProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var membership models.Membership
	if err := initializers.DB.Preload("Project").First(&membership, "user_id=? AND project_id = ?", loggedInUserID, projectID).Error; err != nil {
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
		"message": "You left the project.",
	})
}

func ChangeMemberRole(c *fiber.Ctx) error {
	membershipID := c.Params("membershipID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var reqBody struct {
		Role models.ProjectRole
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var membership models.Membership
	if err := initializers.DB.Preload("Project").First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if membership.Project.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
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
