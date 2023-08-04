package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var messages []models.Message
	if err := paginatedDB.
		Preload("Chat").
		Preload("User").
		Preload("Post").
		Preload("Project").
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if len(messages) > 0 {
		if messages[0].Chat.AcceptingUserID.String() != loggedInUserID && messages[0].Chat.CreatingUserID.String() != loggedInUserID {
			return &fiber.Error{Code: 403, Message: "Cannot perform this action."}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"messages": messages,
	})
}

func GetProjectChatMessages(c *fiber.Ctx) error {
	chatID := c.Params("projectChatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	paginatedDB := API.Paginator(c)(initializers.DB)

	var memberships []models.ProjectChatMembership
	if err := initializers.DB.Where("project_chat_id=?", chatID).Find(&memberships).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	check := false

	for _, membership := range memberships {
		if membership.UserID == parsedLoggedInUserID {
			check = true
		}
	}

	if !check {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	var messages []models.ProjectChatMessage
	if err := paginatedDB.Preload("User").Where("project_chat_id = ? ", chatID).Order("created_at DESC").Find(&messages).Error; err != nil {
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
		Content   string `json:"content"`
		ChatID    string `json:"chatID"`
		PostID    string `json:"postID"`
		ProjectID string `json:"projectID"`
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

func AddProjectChatMessage(c *fiber.Ctx) error {
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

	// var memberships []models.ProjectChatMembership
	// if err := initializers.DB.Where("project_chat_id=?", chatID).Find(&memberships).Error; err != nil {
	// 	return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	// }
	// check := false

	// for _, membership := range memberships {
	// 	if membership.UserID == parsedUserID {
	// 		check = true
	// 	}
	// }

	// if !check {
	// 	return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	// }

	message := models.ProjectChatMessage{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	var chat models.ProjectChat
	if err := initializers.DB.First(&chat, "id=?", parsedChatID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}
	message.ProjectChatID = parsedChatID
	message.ProjectID = chat.ProjectID

	result := initializers.DB.Create(&message)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if err := initializers.DB.Preload("User").First(&message).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

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
