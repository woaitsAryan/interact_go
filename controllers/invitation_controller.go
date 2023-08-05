package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	result = initializers.DB.Save(&invitation)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

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
