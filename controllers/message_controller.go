package controllers

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

func GetMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	// paginatedDB := API.Paginator(c)(initializers.DB)

	var messages []models.Message
	// if err := paginatedDB.
	if err := initializers.DB.
		Preload("Chat").
		Preload("User").
		Preload("Post").
		Preload("Profile").
		Preload("Opening").
		Preload("Opening.Project").
		Preload("Post.User").
		Preload("Project").
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if len(messages) > 0 {
		parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)
		if messages[0].Chat.AcceptingUserID != parsedLoggedInUserID && messages[0].Chat.CreatingUserID != parsedLoggedInUserID {
			return &fiber.Error{Code: 403, Message: "Cannot perform this action."}
		}
		go routines.UpdateChatLastRead(messages[0].ChatID, parsedLoggedInUserID)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"messages": messages,
	})
}

func GetGroupChatMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var membership models.GroupChatMembership
	if err := initializers.DB.Where("group_chat_id=? AND user_id = ?", chatID, loggedInUserID).First(&membership).Error; err != nil {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	var messages []models.GroupChatMessage
	if err := paginatedDB.Preload("User").Where("chat_id = ? ", chatID).Order("created_at DESC").Find(&messages).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		Content string `json:"content"`
		ChatID  string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	var chat models.Chat
	if err := initializers.DB.First(&chat, "id=? AND (creating_user_id = ? OR accepting_user_id = ?)", parsedChatID, parsedUserID, parsedUserID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	message := models.Message{
		UserID:  parsedUserID,
		Content: reqBody.Content,
		ChatID:  parsedChatID,
	}

	result := initializers.DB.Create(&message)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}

func AddGroupChatMessage(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content string `json:"content"`
		ChatID  string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	//TODO add role checks

	chatID := reqBody.ChatID

	var membership models.GroupChatMembership
	if err := initializers.DB.Preload("User").Preload("GroupChat").Where("chat_id=? AND user_id = ?", chatID, loggedInUserID).First(&membership).Error; err != nil {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	if membership.GroupChat.AdminOnly && membership.Role == models.ChatMember {
		return &fiber.Error{Code: 403, Message: "Only admins can send message in this chat."}
	}

	message := models.GroupChatMessage{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	message.ChatID = parsedChatID

	result := initializers.DB.Create(&message)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	message.User = membership.User

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}

func DeleteMessage(c *fiber.Ctx) error {
	messageID := c.Params("messageID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedMessageID, err := uuid.Parse(messageID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var message models.Message
	if err := initializers.DB.First(&message, "id = ? AND user_id=?", parsedMessageID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Message of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	result := initializers.DB.Delete(&message)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Message Deleted",
	})
}
