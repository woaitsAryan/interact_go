package organization_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetNonMembers(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.Where("id = ?", orgID).Preload("Memberships").First(&organization).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	var membershipUserIDs []string

	for _, membership := range organization.Memberships {
		membershipUserIDs = append(membershipUserIDs, membership.UserID.String())
	}

	membershipUserIDs = append(membershipUserIDs, organization.UserID.String())

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

func GetMemberships(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.Where("id = ?", orgID).
		Preload("Memberships").
		Preload("Memberships.User").
		Preload("Invitations").
		Preload("Invitations.User").
		First(&organization).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"organization": organization,
	})
}

func AddMember(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	orgID := c.Params("orgID")

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}

	parsedOrganizationID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	var reqBody struct {
		UserID string `json:"userID"`
		Title  string `json:"title"`
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
			return &fiber.Error{Code: 400, Message: "No Organization of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.UserID == organization.UserID.String() {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.Where("user_id=? AND organization_id=?", user.ID, parsedOrganizationID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var existingInvitation models.Invitation
			err := initializers.DB.Where("user_id=? AND organization_id=? AND status=0", user.ID, parsedOrganizationID).First(&existingInvitation).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "Have already invited this User."}
			}

			var invitation models.Invitation
			invitation.OrganizationID = &parsedOrganizationID
			invitation.UserID = user.ID
			invitation.Title = reqBody.Title

			result := initializers.DB.Create(&invitation)

			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			invitation.User = user
			
			go routines.MarkOrganizationHistory(parsedOrganizationID, parsedUserID, 3, nil, nil, nil, nil, &invitation.ID)

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

	orgMemberID := c.GetRespHeader("orgMemberID")
	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
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

	if membership.UserID.String() == orgMemberID {
		return &fiber.Error{Code: 400, Message: "Cannot remove yourself using this route."}
	}

	if membership.Organization.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
	}

	result := initializers.DB.Delete(&membership)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	go routines.MarkOrganizationHistory(membership.OrganizationID, parsedOrgMemberID, 5, nil, nil, nil, nil, nil )

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User removed to the project.",
	})
}

func LeaveOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var membership models.OrganizationMembership
	if err := initializers.DB.Preload("Organization").First(&membership, "user_id=? AND organization_id = ?", loggedInUserID, orgID).Error; err != nil {
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

	orgChangedUserID := c.GetRespHeader("loggedInUserID")
	loggedInUserID := c.GetRespHeader("orgMemberID")

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var reqBody struct {
		Role  models.OrganizationRole `json:"role"`
		Title string                  `json:"title"`
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

	membership.Title = reqBody.Title

	if orgChangedUserID == loggedInUserID {
		membership.Role = reqBody.Role
	} else {
		if reqBody.Role != models.Manager && membership.Role != models.Manager {
			membership.Role = reqBody.Role
		} else {
			return &fiber.Error{Code: 403, Message: "You don't have the privileges to perform this action."}
		}
	}

	result := initializers.DB.Save(&membership)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User membership updated.",
	})
}
