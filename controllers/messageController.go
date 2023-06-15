package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetChatMessages(c *fiber.Ctx) error {
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

func GetProjectMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var messages []models.ProjectChatMessage
	if err := paginatedDB.Preload("User").Where("chat_id = ? ", chatID).Order("created_at DESC").Find(&messages).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Messages."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"messages": messages,
	})
}

func AddChatMessage(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var chat models.Chat
	if err := initializers.DB.First(&chat, "id=?", parsedChatID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	var reqBody struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
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

func AddProjectMessage(c *fiber.Ctx) error {
	chatID := c.Params("projectChatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedProjectChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var chat models.ProjectChatMessage
	if err := initializers.DB.First(&chat, "id=?", parsedProjectChatID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	var reqBody struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	message := models.ProjectChatMessage{
		UserID:        parsedUserID,
		ProjectChatID: parsedProjectChatID,
		Content:       reqBody.Content,
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

func DeleteChatMessage(c *fiber.Ctx) error {
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

func DeleteProjectMessage(c *fiber.Ctx) error {
	messageID := c.Params("messageID")

	parsedMessageID, err := uuid.Parse(messageID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var message models.ProjectChatMessage
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
