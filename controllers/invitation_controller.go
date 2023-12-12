package controllers

import (
	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetInvitations(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var invitations []models.Invitation
	if err := initializers.DB.Preload("GroupChat").
		Preload("Project").
		Preload("Organization").
		Preload("Organization.User").
		Where("user_id = ? ", loggedInUserID).Order("created_at DESC").Find(&invitations).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "",
		"invitations": invitations,
	})
}

func AcceptInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("invitationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var user models.User
	if err := initializers.DB.Where("id=?", loggedInUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	if !user.Verified {
		return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
	}

	var invitation models.Invitation
	err := initializers.DB.Preload("Project").First(&invitation, "id=? AND user_id=?", invitationID, loggedInUserID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	if invitation.Status != 0 {
		return &fiber.Error{Code: 400, Message: "Cannot Perform this action."}
	}

	invitation.Status = 1

	if invitation.ProjectID != nil {
		membership := models.Membership{
			UserID:    invitation.UserID,
			ProjectID: *invitation.ProjectID,
			Title:     invitation.Title,
			Role:      models.ProjectMember,
		}

		result := initializers.DB.Create(&membership)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
		}

		go routines.MarkProjectHistory(*invitation.ProjectID, parsedLoggedInUserID, 1, nil, nil, nil, nil, nil, "")
		cache.RemoveProject(invitation.Project.Slug)
		cache.RemoveProject("-workspace--" + invitation.Project.Slug)
	} else if invitation.OrganizationID != nil {
		membership := models.OrganizationMembership{
			UserID:         invitation.UserID,
			OrganizationID: *invitation.OrganizationID,
			Role:           models.Member,
			Title:          invitation.Title,
		}

		result := initializers.DB.Create(&membership)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
		}
	} else if invitation.GroupChatID != nil {
		membership := models.GroupChatMembership{
			UserID:      invitation.UserID,
			GroupChatID: *invitation.GroupChatID,
			Role:        models.ChatMember,
		}

		result := initializers.DB.Create(&membership)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
		}
	}

	result := initializers.DB.Save(&invitation)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}
	if(invitation.OrganizationID != nil){
		go routines.IncrementOrgMember(*invitation.OrganizationID)
	}
	if(invitation.ProjectID != nil){
		go routines.IncrementProjectMember(*invitation.ProjectID)
	}

	go routines.SendInvitationAcceptedNotification(invitation.UserID, parsedLoggedInUserID)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func RejectInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("invitationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var invitation models.Invitation
	err := initializers.DB.First(&invitation, "id=? AND user_id=?", invitationID, loggedInUserID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	if invitation.Status != 0 {
		return &fiber.Error{Code: 400, Message: "Cannot Perform this action."}
	}

	invitation.Status = -1

	result := initializers.DB.Save(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func WithdrawInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("invitationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var invitation models.Invitation
	err := initializers.DB.Preload("Project").Preload("Organization").First(&invitation, "id=?", invitationID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}
	if invitation.ProjectID != nil {
		if invitation.Project.UserID.String() != loggedInUserID {
			return &fiber.Error{Code: 403, Message: "You don't have the permission to perform this action."}
		}
		cache.RemoveProject("-workspace--" + invitation.Project.Slug)
	} else if invitation.OrganizationID != nil {
		if invitation.Organization.UserID.String() != loggedInUserID {
			return &fiber.Error{Code: 403, Message: "You don't have the permission to perform this action."}
		}
	}

	if invitation.Status == 1 {
		return &fiber.Error{Code: 400, Message: "Invitation is already accepted, cannot withdraw now."}
	}

	result := initializers.DB.Delete(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")
	if orgMemberID != "" && orgID != ""{
		parsedOrgID, err := uuid.Parse(orgID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
		}

		parsedOrgMemberID, err := uuid.Parse(orgMemberID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID."}
		}
		go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 4, nil, nil, nil, nil, nil, invitation.Title) 
	}


	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Withdrawn",
	})
}

func GetUnreadInvitationCount(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var count int64
	if err := initializers.DB.
		Model(models.Invitation{}).
		Where("user_id=? AND status=0", loggedInUserID).
		Count(&count).
		Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"count":   count,
	})
}

func MarkReadInvitations(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		UnreadInvitations []string `json:"unreadInvitations"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	for _, unreadInvitationID := range reqBody.UnreadInvitations {
		var invitation models.Invitation
		if err := initializers.DB.
			Where("id=? AND user_id=?", unreadInvitationID, loggedInUserID).
			First(&invitation).
			Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
		invitation.Read = true
		result := initializers.DB.Save(&invitation)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
	})
}
