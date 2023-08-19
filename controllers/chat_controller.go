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

func GetChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chat models.Chat

	if err := initializers.DB.Preload("Messages").First(&chat, "id=? AND (creating_user_id = ? OR accepting_user_id = ?)", chatID, loggedInUserID, loggedInUserID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chat":    chat,
	})
}

func GetUserNonPopulatedChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chats []models.Chat
	if err := initializers.DB.Where("creating_user_id=?", loggedInUserID).Or("accepting_user_id = ?", loggedInUserID).Find(&chats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chats":   chats,
	})
}

// TODO separate get personal chats and group chats
func GetUserChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chats []models.Chat
	if err := initializers.DB.
		Preload("CreatingUser").
		Preload("AcceptingUser").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Where("creating_user_id=?", loggedInUserID).
		Or("accepting_user_id = ?", loggedInUserID).
		Find(&chats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var groupChats []models.GroupChat
	if err := initializers.DB.
		Preload("Project").
		Preload("Organization").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Preload("Memberships").
		Preload("Memberships.User").
		Joins("JOIN group_chat_memberships ON group_chat_memberships.group_chat_id = group_chats.id").
		Where("group_chat_memberships.user_id = ?", loggedInUserID).
		Find(&groupChats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "",
		"chats":      chats,
		"groupChats": groupChats,
	})
}

func AcceptChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chat models.Chat
	if err := initializers.DB.First(&chat, "id = ? AND accepting_user_id=?", chatID, loggedInUserID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No chat of this id found."}
	}

	chat.Accepted = true

	result := initializers.DB.Save(&chat)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Accepted",
		"chat":    chat,
	})
}

func AddChat(c *fiber.Ctx) error {
	var reqBody struct {
		UserID string `json:"userID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	userID := c.GetRespHeader("loggedInUserID")
	chatUserID := reqBody.UserID

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}
	parsedChatUserID, err := uuid.Parse(chatUserID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid User ID."}
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", parsedUserID)
	var chatUser models.User
	if err = initializers.DB.First(&chatUser, "id = ?", parsedChatUserID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No user of this id found."}
	}

	var existingChat models.Chat
	if err = initializers.DB.Where("creating_user_id = ? AND accepting_user_id = ?", parsedUserID, parsedChatUserID).
		Or("creating_user_id = ? AND accepting_user_id = ?", parsedChatUserID, parsedUserID).
		First(&existingChat).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "Chat already exists between the users."}
	}

	chat := models.Chat{
		CreatingUserID:  parsedUserID,
		AcceptingUserID: chatUser.ID,
	}

	result := initializers.DB.Create(&chat)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	go routines.SendChatNotification(parsedUserID, parsedChatUserID)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Created",
		"chat":    chat,
	})
}

func DeleteChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var chat models.Chat
	if err := initializers.DB.First(&chat, "id = ?", parsedChatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if err := initializers.DB.Delete(&chat).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully",
	})
}
