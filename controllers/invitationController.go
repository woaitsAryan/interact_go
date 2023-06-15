package controllers

import (
	"sort"
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func sortInvitationsByCreatedAt(invitations []interface{}) {
	sort.Slice(invitations, func(i, j int) bool {
		invitationA := invitations[i]
		invitationB := invitations[j]

		createdAtA := getCreatedAt(invitationA)
		createdAtB := getCreatedAt(invitationB)

		return createdAtA.Before(createdAtB)
	})
}

func getCreatedAt(invitation interface{}) time.Time {
	switch inv := invitation.(type) {
	case models.ChatInvitation:
		return inv.CreatedAt
	case models.ProjectInvitation:
		return inv.CreatedAt
	default:
		return time.Time{}
	}
}

func GetInvitations(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chatInvitations []models.ChatInvitation
	if err := initializers.DB.Preload("Chat").Where("user_id = ? ", loggedInUserID).Order("created_at DESC").Find(&chatInvitations).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Chat Invitations."}
	}

	var projectInvitations []models.ProjectInvitation
	if err := initializers.DB.Preload("Project").Where("user_id = ? ", loggedInUserID).Order("created_at DESC").Find(&projectInvitations).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Project Invitations."}
	}

	var combinedInvitations []interface{}

	for _, chatInvitation := range chatInvitations {
		combinedInvitations = append(combinedInvitations, chatInvitation)
	}

	for _, projectInvitation := range projectInvitations {
		combinedInvitations = append(combinedInvitations, projectInvitation)
	}

	sortInvitationsByCreatedAt(combinedInvitations)

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "",
		"invitations": combinedInvitations,
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
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the invitation."}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func AcceptProjectInvitation(c *fiber.Ctx) error {

	invitationID := c.Params("invitationID")

	var invitation models.ProjectInvitation
	err := initializers.DB.First(&invitation, "id=?", invitationID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	invitation.Status = 1

	result := initializers.DB.Save(&invitation)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the invitation."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func RejectChatInvitation(c *fiber.Ctx) error {

	invitationID := c.Params("invitationID")

	var invitation models.ChatInvitation
	err := initializers.DB.First(&invitation, "id=?", invitationID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	invitation.Status = -1

	result := initializers.DB.Save(&invitation)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the invitation."}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Accepted",
	})
}

func RejectProjectInvitation(c *fiber.Ctx) error {

	invitationID := c.Params("invitationID")

	var invitation models.ProjectInvitation
	err := initializers.DB.First(&invitation, "id=?", invitationID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	invitation.Status = -1

	result := initializers.DB.Save(&invitation)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating the invitation."}
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
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the invitation."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Withdrawn",
	})
}

func WithdrawProjectInvitation(c *fiber.Ctx) error {

	invitationID := c.Params("invitationID")

	var invitation models.ProjectInvitation
	err := initializers.DB.First(&invitation, "id=?", invitationID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Invitation of this ID found."}
	}

	result := initializers.DB.Delete(&invitation)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the invitation."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitation Withdrawn",
	})
}
