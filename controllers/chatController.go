package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	var chat models.Chat

	err := initializers.DB.Preload("Messages").First(&chat, "id=?", chatID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chat":    chat,
	})
}

func GetProjectChat(c *fiber.Ctx) error {
	chatID := c.Params("projectChatID")

	var chat models.ProjectChat

	err := initializers.DB.First(&chat, "id=?", chatID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chat":    chat,
	})
}

func GetUserNonPopulatedChats(c *fiber.Ctx) error { // ! GET USER CHATS
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chats []models.Chat
	if err := initializers.DB.Where("creating_user_id=?", loggedInUserID).Or("accepting_user_id = ?", loggedInUserID).Find(&chats).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chats":   chats,
	})
}

func GetUserChats(c *fiber.Ctx) error { // ! GET USER CHATS
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chats []models.Chat
	if err := initializers.DB.Preload("CreatingUser").Preload("AcceptingUser").Where("creating_user_id=?", loggedInUserID).Or("accepting_user_id = ?", loggedInUserID).Find(&chats).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	// var projectChats models.ProjectChat
	// if err := initializers.DB.Where("user_id=?", loggedInUserID).Find(&projectChats).Error; err != nil {
	// 	return &fiber.Error{Code: 500, Message: "Database Error."}
	// }

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chats":   chats,
		// "projectChats": projectChats,
	})
}

func AcceptChat(c *fiber.Ctx) error {

	chatID := c.Params("chatID")

	var chat models.Chat
	err := initializers.DB.First(&chat, "id = ?", chatID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No chat of this id found."}
	}

	chat.Accepted = true

	result := initializers.DB.Save(&chat)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating chat."}
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
	err = initializers.DB.First(&chatUser, "id = ?", parsedChatUserID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No user of this id found."}
	}

	var existingChat models.Chat
	err = initializers.DB.Where("creating_user_id = ? AND accepting_user_id = ?", parsedUserID, parsedChatUserID).
		Or("creating_user_id = ? AND accepting_user_id = ?", parsedChatUserID, parsedUserID).
		First(&existingChat).Error

	if err == nil {
		return &fiber.Error{Code: 400, Message: "Chat already exists between the users."}
	}

	chat := models.Chat{
		CreatingUserID:  parsedUserID,
		AcceptingUserID: chatUser.ID,
	}

	result := initializers.DB.Create(&chat)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating chat"}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Created",
		"chat":    chat,
	})
}

func AddGroupChat(c *fiber.Ctx) error { //! ADD VALIDATORS
	var reqBody struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		UserIDs     []string `json:"userIDs"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	userID := c.GetRespHeader("loggedInUserID")
	chatUserIDs := reqBody.UserIDs

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	chatUsers := make([]models.User, len(chatUserIDs))

	for i, chatUserID := range chatUserIDs {
		parsedChatUserID, err := uuid.Parse(chatUserID)
		if err != nil {
			return &fiber.Error{Code: 500, Message: "Invalid User ID."}
		}
		var chatUser models.User
		err = initializers.DB.First(&chatUser, "creating_user_id = ?", parsedChatUserID).Error

		if err != nil {
			return &fiber.Error{Code: 400, Message: "No user of this id found."}
		}

		chatUsers[i] = chatUser
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", parsedUserID)

	chat := models.GroupChat{
		CreatingUserID: parsedUserID,
		Title:          reqBody.Title,
		Description:    reqBody.Description,
		Members:        []models.User{user},
	}

	result := initializers.DB.Create(&chat)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating chat"}
	}

	// for _, chatUser := range chatUsers {
	// 	invitation := models.ChatInvitation{
	// 		UserID: chatUser.ID,
	// 		ChatID: chat.ID,
	// 	}
	// 	result := initializers.DB.Create(&invitation)

	// 	if result.Error != nil {
	// 		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating invitations"}
	// 	}
	// }

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Created",
		"chat":    chat,
	})
}

func AddProjectChat(c *fiber.Ctx) error {
	var reqBody struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		UserIDs     []string `json:"userIDs"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	userID := c.GetRespHeader("loggedInUserID")
	chatUserIDs := reqBody.UserIDs

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid Project ID."}
	}

	chatUsers := make([]models.User, len(chatUserIDs)+1)

	for i, chatUserID := range chatUserIDs {
		parsedChatUserID, err := uuid.Parse(chatUserID)
		if err != nil {
			return &fiber.Error{Code: 500, Message: "Invalid User ID."}
		}
		var chatUser models.User
		err = initializers.DB.First(&chatUser, "creating_user_id = ?", parsedChatUserID).Error

		if err != nil {
			return &fiber.Error{Code: 400, Message: "No user of this id found."}
		}

		chatUsers[i] = chatUser
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", parsedUserID)

	chatUsers = append(chatUsers, user)

	var project models.Project
	err = initializers.DB.First(&project, "id = ?", parsedProjectID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No project of this id found."}
	}

	projectChat := models.ProjectChat{
		CreatingUserID: parsedUserID,
		Title:          reqBody.Title,
		Description:    reqBody.Description,
		ProjectID:      project.ID,
		Members:        chatUsers,
	}

	result := initializers.DB.Create(&projectChat)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating chat"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Created",
		"chat":    projectChat,
	})
}

func EditGroupChat(c *fiber.Ctx) error { //* Adding new users here only
	var reqBody struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		UserIDs     []string `json:"userIDs"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	newChatUserIDs := reqBody.UserIDs
	newChatUsers := make([]models.User, len(newChatUserIDs))

	for i, newChatUserID := range newChatUserIDs {
		parsedNewChatUserID, err := uuid.Parse(newChatUserID)
		if err != nil {
			return &fiber.Error{Code: 500, Message: "Invalid User ID."}
		}
		var chatUser models.User
		err = initializers.DB.First(&chatUser, "id = ?", parsedNewChatUserID).Error

		if err != nil {
			return &fiber.Error{Code: 400, Message: "No user of this id found."}
		}

		newChatUsers[i] = chatUser
	}

	chatID := c.Params("chatID")

	var chat models.GroupChat
	err := initializers.DB.First(&chat, "id = ?", chatID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No chat of this id found."}
	}

	if reqBody.Title != "" {
		chat.Title = reqBody.Title
	}
	if reqBody.Description != "" {
		chat.Description = reqBody.Description
	}

	result := initializers.DB.Save(&chat)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating chat."}
	}

	for _, chatUser := range newChatUsers {
		invitation := models.ChatInvitation{
			UserID:      chatUser.ID,
			GroupChatID: chat.ID,
		}
		result := initializers.DB.Create(&invitation)

		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while creating invitations"}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Updated",
		"chat":    chat,
	})
}

func EditProjectChat(c *fiber.Ctx) error { //* Adding new users here only
	var reqBody struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		UserIDs     []string `json:"userIDs"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	newChatUserIDs := reqBody.UserIDs
	newChatUsers := make([]models.User, len(newChatUserIDs))

	for i, newChatUserID := range newChatUserIDs {
		parsedNewChatUserID, err := uuid.Parse(newChatUserID)
		if err != nil {
			return &fiber.Error{Code: 500, Message: "Invalid User ID."}
		}
		var chatUser models.User
		err = initializers.DB.First(&chatUser, "id = ?", parsedNewChatUserID).Error

		if err != nil {
			return &fiber.Error{Code: 400, Message: "No user of this id found."}
		}

		newChatUsers[i] = chatUser
	}

	projectChatID := c.Params("projectChatID")

	var projectChat models.ProjectChat
	err := initializers.DB.First(&projectChat, "id = ?", projectChatID).Error

	if err != nil {
		return &fiber.Error{Code: 400, Message: "No chat of this id found."}
	}

	if reqBody.Title != "" {
		projectChat.Title = reqBody.Title
	}
	if reqBody.Description != "" {
		projectChat.Description = reqBody.Description
	}
	if reqBody.UserIDs != nil {
		projectChat.Members = append(projectChat.Members, newChatUsers...)
	}

	result := initializers.DB.Save(&projectChat)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while updating chat."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Updated",
		"chat":    projectChat,
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
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&chat).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully",
	})
}

func DeleteGroupChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var chat models.GroupChat
	if err := initializers.DB.First(&chat, "id = ?", parsedChatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&chat).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully",
	})
}

func DeleteProjectChat(c *fiber.Ctx) error {
	chatID := c.Params("projectChatID")

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var chat models.ProjectChat
	if err := initializers.DB.First(&chat, "id = ?", parsedChatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&chat).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully",
	})
}
