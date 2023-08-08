package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// func sortInvitationsByCreatedAt(invitations []interface{}) {
// 	sort.Slice(invitations, func(i, j int) bool {
// 		invitationA := invitations[i]
// 		invitationB := invitations[j]

// 		createdAtA := getCreatedAt(invitationA)
// 		createdAtB := getCreatedAt(invitationB)

// 		return createdAtA.Before(createdAtB)
// 	})
// }

// func getCreatedAt(invitation interface{}) time.Time {
// 	switch inv := invitation.(type) {
// 	case models.ChatInvitation:
// 		return inv.CreatedAt
// 	case models.ProjectInvitation:
// 		return inv.CreatedAt
// 	default:
// 		return time.Time{}
// 	}
// }

func GetInvitations(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	// var chatInvitations []models.ChatInvitation
	// if err := initializers.DB.Preload("Chat").Where("user_id = ? ", loggedInUserID).Order("created_at DESC").Find(&chatInvitations).Error; err != nil {
	// 	return &fiber.Error{Code: 500, Message: "Failed to get the Chat Invitations."}
	// }

	var projectInvitations []models.ProjectInvitation
	if err := initializers.DB.Preload("Project").Where("user_id = ? ", loggedInUserID).Order("created_at DESC").Find(&projectInvitations).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	// var combinedInvitations []interface{}

	// for _, chatInvitation := range chatInvitations {
	// 	combinedInvitations = append(combinedInvitations, chatInvitation)
	// }

	// for _, projectInvitation := range projectInvitations {
	// 	combinedInvitations = append(combinedInvitations, projectInvitation)
	// }

	// sortInvitationsByCreatedAt(combinedInvitations)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		// "chatInvitations":    chatInvitations,
		"projectInvitations": projectInvitations,
	})
}

func AcceptChatInvitation(c *fiber.Ctx) error {

	invitationID := c.Params("invitationID")

	var invitation models.ChatInvitation
	err := initializers.DB.First(&invitation, "id=?", invitationID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	invitation.Status = 1

	result := initializers.DB.Save(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func AcceptProjectInvitation(c *fiber.Ctx) error {
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

	var invitation models.ProjectInvitation
	err := initializers.DB.First(&invitation, "id=? AND user_id=?", invitationID, loggedInUserID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	if invitation.Status != 0 {
		return &fiber.Error{Code: 400, Message: "Cannot Perform this action."}
	}

	invitation.Status = 1

	membership := models.Membership{
		UserID:    invitation.UserID,
		ProjectID: invitation.ProjectID,
		Title:     invitation.Title,
		Role:      "Member",
	}

	result := initializers.DB.Create(&membership)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	result = initializers.DB.Save(&invitation)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	go routines.SendInvitationAcceptedNotification(invitation.UserID, parsedLoggedInUserID)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func RejectChatInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("invitationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var invitation models.ChatInvitation
	err := initializers.DB.First(&invitation, "id=? AND user_id=?", invitationID, loggedInUserID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	invitation.Status = -1

	result := initializers.DB.Save(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func RejectProjectInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("invitationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var invitation models.ProjectInvitation
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

func WithdrawChatInvitation(c *fiber.Ctx) error {

	invitationID := c.Params("invitationID")

	var invitation models.ChatInvitation
	err := initializers.DB.First(&invitation, "id=?", invitationID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	result := initializers.DB.Delete(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Withdrawn",
	})
}

func WithdrawProjectInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("invitationID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var invitation models.ProjectInvitation
	err := initializers.DB.Preload("Project").First(&invitation, "id=?", invitationID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	if invitation.Project.UserID.String() != loggedInUserID {
		return &fiber.Error{Code: 403, Message: "You don't have the permission to perform this action."}
	}

	if invitation.Status == 1 {
		return &fiber.Error{Code: 400, Message: "Invitation is already accepted, cannot withdraw now."}
	}

	result := initializers.DB.Delete(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Withdrawn",
	})
}

func GetUnreadInvitations(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var count int64
	if err := initializers.DB.
		Model(models.ProjectInvitation{}).
		Where("user_id=? AND read=?", loggedInUserID, false).
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
		var invitation models.ProjectInvitation
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
