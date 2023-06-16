package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var messages []models.Message
	if err := paginatedDB.Preload("User").Where("chat_id = ? ", chatID).Order("created_at DESC").Find(&messages).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Messages."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"messages": messages,
	})
}

func AddMessage(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content       string `json:"content"`
		ChatID        string `json:"chatID"`
		ProjectChatID string `json:"projectChatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID
	projectChatID := reqBody.ProjectChatID

	var parsedChatID uuid.UUID

	if chatID != "" {
		parsedChatID, err := uuid.Parse(chatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		var chat models.Chat
		if err := initializers.DB.First(&chat, "id=?", parsedChatID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
		}
	} else if projectChatID != "" {
		parsedChatID, err := uuid.Parse(projectChatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		var chat models.ProjectChat
		if err := initializers.DB.First(&chat, "id=?", parsedChatID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
		}
	} else {
		return &fiber.Error{Code: 400, Message: "Invalid Chat ID."}
	}

	message := models.Message{
		UserID:  parsedUserID,
		ChatID:  parsedChatID,
		Content: reqBody.Content,
	}

	result := initializers.DB.Create(&message)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the message."}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Message Added",
	})
}

func DeleteMessage(c *fiber.Ctx) error {
	messageID := c.Params("messageID")

	parsedMessageID, err := uuid.Parse(messageID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var message models.Message
	if err := initializers.DB.First(&message, "id = ?", parsedMessageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Message of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	result := initializers.DB.Delete(&message)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while deleting the message."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Message Deleted",
	})
}
