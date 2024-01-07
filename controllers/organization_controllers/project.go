package organization_controllers

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddProjectMembers(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	var project models.Project
	if err := initializers.DB.Preload("User").Preload("Memberships").First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if project.Memberships != nil && len(project.Memberships) > 0 {
		return &fiber.Error{Code: 403, Message: "Project already has memberships, cannot perform this action."}
	}

	type UserSlice struct {
		UserID string             `json:"userID"`
		Title  string             `json:"title"`
		Role   models.ProjectRole `json:"role"`
	}

	type ReqBody struct {
		UserSlices []UserSlice `json:"userSlices"`
	}

	var reqBody ReqBody
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	orgID := c.Params("orgID")

	var memberships []models.Membership

	for _, slice := range reqBody.UserSlices {
		var orgMembership models.OrganizationMembership
		if err := initializers.DB.First(&orgMembership, "user_id = ? AND organization_id=?", slice.UserID, orgID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		membership := models.Membership{
			ProjectID: project.ID,
			UserID:    orgMembership.UserID,
			Title:     slice.Title,
			Role:      slice.Role,
		}

		result := initializers.DB.Create(&membership)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
		}

		memberships = append(memberships, membership)
	}

	project.NumberOfMembers = len(memberships)

	result := initializers.DB.Save(&project)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	project.Memberships = memberships

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Members Added",
		"project": project,
	})
}
