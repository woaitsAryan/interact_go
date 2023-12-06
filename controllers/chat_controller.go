package controllers

import (
	"time"

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
	if err := initializers.DB.
		Preload("CreatingUser").
		Preload("AcceptingUser").
		Preload("LatestMessage").
		Preload("LastReadMessageByAcceptingUser").
		Preload("LastReadMessageByCreatingUser").
		First(&chat, "id=? AND (creating_user_id = ? OR accepting_user_id = ?)", chatID, loggedInUserID, loggedInUserID).Error; err != nil {
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

	var groupChats []models.GroupChat
	if err := initializers.DB.
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

func GetPersonalUnFilteredChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chats []models.Chat
	if err := initializers.DB.
		Preload("CreatingUser").
		Preload("AcceptingUser").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Where("creating_user_id=? OR accepting_user_id = ?", loggedInUserID, loggedInUserID).
		Find(&chats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chats":   chats,
	})
}

func GetPersonalChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var chats []models.Chat
	if err := initializers.DB.
		Preload("CreatingUser").
		Preload("AcceptingUser").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Where("creating_user_id=? OR accepting_user_id = ?", loggedInUserID, loggedInUserID).
		Find(&chats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var filteredChats []models.Chat
	var requests []models.Chat

	for _, chat := range chats {
		if chat.Accepted {
			if chat.CreatingUserID.String() == loggedInUserID {
				if chat.LatestMessage != nil {
					if chat.LatestMessage.CreatedAt.After(chat.LastResetByCreatingUser) {
						filteredChats = append(filteredChats, chat)
					}
				}
			} else {
				if chat.LatestMessage != nil {
					if chat.LatestMessage.CreatedAt.After(chat.LastResetByAcceptingUser) {
						filteredChats = append(filteredChats, chat)
					}
				}
			}
		} else {
			if chat.CreatingUserID.String() == loggedInUserID {
				if chat.LatestMessage != nil {
					filteredChats = append(filteredChats, chat)
				}
			} else {
				requests = append(requests, chat)
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"chats":    filteredChats,
		"requests": requests,
	})
}

func GetGroupChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var groupChats []models.GroupChat
	if err := initializers.DB.
		Preload("User").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Preload("Memberships").
		Preload("Memberships.User").
		Joins("JOIN group_chat_memberships ON group_chat_memberships.group_chat_id = group_chats.id").
		Where("group_chat_memberships.user_id = ? AND group_chats.project_id IS NULL", loggedInUserID).
		Find(&groupChats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chats":   groupChats,
	})
}

func GetProjectChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var groupChats []models.GroupChat
	if err := initializers.DB.
		Preload("User").
		Preload("Project").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Preload("Memberships").
		Preload("Memberships.User").
		Joins("JOIN group_chat_memberships ON group_chat_memberships.group_chat_id = group_chats.id").
		Where("group_chat_memberships.user_id = ? AND group_chats.project_id IS NOT NULL", loggedInUserID).
		Find(&groupChats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var projects []models.Project

	for _, chat := range groupChats {
		check := false
		for _, project := range projects {
			if *chat.ProjectID == project.ID {
				check = true
				break
			}
		}
		if !check {
			projects = append(projects, chat.Project)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"chats":    groupChats,
		"projects": projects,
	})
}

func GetOrgChats(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var groupChats []models.GroupChat
	if err := initializers.DB.
		Preload("User").
		Preload("Organization").
		Preload("LatestMessage").
		Preload("LatestMessage.User").
		Preload("Memberships").
		Preload("Memberships.User").
		Joins("JOIN group_chat_memberships ON group_chat_memberships.group_chat_id = group_chats.id").
		Where("group_chat_memberships.user_id = ? AND group_chats.organization_id IS NOT NULL", loggedInUserID).
		Find(&groupChats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var organizations []models.Organization

	for _, chat := range groupChats {
		check := false
		for _, project := range organizations {
			if *chat.ProjectID == project.ID {
				check = true
				break
			}
		}
		if !check {
			organizations = append(organizations, chat.Organization)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"message":       "",
		"chats":         groupChats,
		"organizations": organizations,
	})
}

func GetUnreadChats(c *fiber.Ctx) error { //* Personal Only
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var chats []models.Chat
	if err := initializers.DB.
		Where("creating_user_id=? OR accepting_user_id = ?", loggedInUserID, loggedInUserID).
		Find(&chats).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var message models.Message
	var chatIDs []string
	for _, chat := range chats {
		if chat.AcceptingUserID == parsedLoggedInUserID {
			if err := initializers.DB.Where("chat_id=? AND user_id=?", chat.ID, chat.CreatingUserID).
				Order("created_at DESC").
				First(&message).Error; err == nil {
				if chat.LastReadMessageByAcceptingUserID != nil && *chat.LastReadMessageByAcceptingUserID != message.ID {
					chatIDs = append(chatIDs, chat.ID.String())
				}
			}
		} else if chat.CreatingUserID == parsedLoggedInUserID {
			if err := initializers.DB.Where("chat_id=? AND user_id=?", chat.ID, chat.AcceptingUserID).
				Order("created_at DESC").
				First(&message).Error; err == nil {
				if chat.LastReadMessageByCreatingUserID != nil && *chat.LastReadMessageByCreatingUserID != message.ID {
					chatIDs = append(chatIDs, chat.ID.String())
				}
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chatIDs": chatIDs,
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
	go helpers.SendChatMail(chatUser.Name, chatUser.Email, user.Name)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Created",
		"chat":    chat,
	})
}

func UpdateLastRead(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)
	chatID := c.Params("chatID")

	var chat models.Chat
	if err := initializers.DB.Where("id = ? AND (creating_user_id = ? OR accepting_user_id = ?)", chatID, parsedLoggedInUserID, parsedLoggedInUserID).
		First(&chat).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	var message models.Message

	if chat.AcceptingUserID == parsedLoggedInUserID {
		if err := initializers.DB.Where("chat_id=? AND user_id=?", chat.ID, chat.CreatingUserID).
			Order("created_at DESC").
			First(&message).Error; err == nil {
			chat.LastReadMessageByAcceptingUserID = &message.ID
		}
	} else if chat.CreatingUserID == parsedLoggedInUserID {
		if err := initializers.DB.Where("chat_id=? AND user_id=?", chat.ID, chat.AcceptingUserID).
			Order("created_at DESC").
			First(&message).Error; err == nil {
			chat.LastReadMessageByCreatingUserID = &message.ID
		}
	}

	result := initializers.DB.Save(&chat)
	if result.Error != nil {
		helpers.LogDatabaseError("Error while updating Chat-UpdateChatLastRead", result.Error, "go_routine")
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Last Read Updated",
	})
}

func BlockChat(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		ChatID string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid Chat ID."}
	}

	var chat models.Chat
	if err = initializers.DB.Where("id = ? AND (creating_user_id = ? OR accepting_user_id = ?)", parsedChatID, parsedLoggedInUserID, parsedLoggedInUserID).
		First(&chat).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	if parsedLoggedInUserID == chat.AcceptingUserID {
		chat.BlockedByAcceptingUser = true
	} else if parsedLoggedInUserID == chat.CreatingUserID {
		chat.BlockedByCreatingUser = true
	}

	if err := initializers.DB.Save(&chat).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Blocked",
	})
}

func UnblockChat(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		ChatID string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid Chat ID."}
	}

	var chat models.Chat
	if err = initializers.DB.Where("id = ? AND (creating_user_id = ? OR accepting_user_id = ?)", parsedChatID, parsedLoggedInUserID, parsedLoggedInUserID).
		First(&chat).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	if parsedLoggedInUserID == chat.AcceptingUserID {
		chat.BlockedByAcceptingUser = false
	} else if parsedLoggedInUserID == chat.CreatingUserID {
		chat.BlockedByCreatingUser = false
	}

	if err := initializers.DB.Save(&chat).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Unblocked",
	})
}

func ResetChat(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		ChatID string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Invalid Chat ID."}
	}

	var chat models.Chat
	if err = initializers.DB.Where("id = ? AND (creating_user_id = ? OR accepting_user_id = ?)", parsedChatID, parsedLoggedInUserID, parsedLoggedInUserID).
		First(&chat).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	if parsedLoggedInUserID == chat.AcceptingUserID {
		chat.LastResetByAcceptingUser = time.Now()
	} else if parsedLoggedInUserID == chat.CreatingUserID {
		chat.LastResetByCreatingUser = time.Now()
	}

	if err := initializers.DB.Save(&chat).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Reset",
	})
}

func DeleteChat(c *fiber.Ctx) error { //TODO have a history for project and organizations
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
