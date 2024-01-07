package messaging_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetGroupChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var chat models.GroupChat
	err := initializers.DB.
		Preload("User").
		Preload("Memberships").
		Preload("Memberships.User").
		Preload("Invitations").
		Preload("Invitations.User").
		Preload("Project").
		Preload("Organization").
		Where("id = ?", chatID).
		First(&chat).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	check := false
	var userMembership models.GroupChatMembership
	for _, membership := range chat.Memberships { // Even Owner has a chat membership
		if membership.UserID == parsedLoggedInUserID {
			userMembership = membership
			check = true
		}
	}

	if chat.Project.UserID == parsedLoggedInUserID {
		check = true
	}

	if !check {
		return &fiber.Error{Code: 403, Message: "Do not have the permission to perform this action."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "",
		"chat":       chat,
		"membership": userMembership,
	})
}

func AddGroupChat(chatType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqBody struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			UserIDs     []string `json:"userIDs"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}
		chatUserIDs := reqBody.UserIDs

		var chat models.GroupChat
		var parsedLoggedInUserID uuid.UUID

		if chatType == "Group" {
			loggedInUserID := c.GetRespHeader("loggedInUserID")
			parsedLoggedInUserID, err := uuid.Parse(loggedInUserID)
			if err != nil {
				return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
			}

			picName, err := utils.UploadImage(c, "coverPic", helpers.ChatClient, 720, 720)
			if err != nil {
				return err
			}

			chat := models.GroupChat{
				UserID:      parsedLoggedInUserID,
				Title:       reqBody.Title,
				Description: reqBody.Description,
			}

			if picName != "" {
				chat.CoverPic = picName
			}

			result := initializers.DB.Create(&chat)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}

			chatMembership := models.GroupChatMembership{
				UserID:      parsedLoggedInUserID,
				GroupChatID: chat.ID,
				Role:        models.ChatAdmin,
			}
			result = initializers.DB.Create(&chatMembership)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}

			for _, chatUserID := range chatUserIDs {
				parsedUserID, err := uuid.Parse(chatUserID)
				if err != nil {
					return &fiber.Error{Code: 500, Message: "Error Parsing the User ID."}
				}
				if parsedUserID == parsedLoggedInUserID {
					continue
				}
				invitation := models.Invitation{
					UserID:      parsedUserID,
					GroupChatID: &chat.ID,
				}
				result := initializers.DB.Create(&invitation)

				if result.Error != nil {
					return &fiber.Error{Code: 500, Message: "Internal Server Error while creating invitations"}
				}
			}
		} else {
			if chatType == "Project" { //TODO check for project and organization memberships of the users in reqBody
				userID := c.GetRespHeader("projectMemberID")
				if userID == "" {
					userID = c.GetRespHeader("loggedInUserID")
				}

				parsedUserID, err := uuid.Parse(userID)
				if err != nil {
					return &fiber.Error{Code: 500, Message: "Error Parsing the LoggedIn User ID."}
				}

				parsedLoggedInUserID = parsedUserID

				projectID := c.Params("projectID")

				var project models.Project
				if err := initializers.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}

				picName, err := utils.UploadImage(c, "coverPic", helpers.ChatClient, 720, 720)
				if err != nil {
					return err
				}

				chat = models.GroupChat{
					UserID:      parsedLoggedInUserID,
					Title:       reqBody.Title,
					Description: reqBody.Description,
					ProjectID:   &project.ID,
				}

				if picName != "" {
					chat.CoverPic = picName
				}

				result := initializers.DB.Create(&chat)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}

				projectMemberID := c.GetRespHeader("projectMemberID")
				parsedID, _ := uuid.Parse(projectMemberID)
				go routines.MarkProjectHistory(project.ID, parsedID, 8, nil, nil, nil, nil, nil, "")

			} else if chatType == "Organization" {
				userID := c.GetRespHeader("orgMemberID")
				parsedUserID, _ := uuid.Parse(userID)

				parsedLoggedInUserID = parsedUserID

				orgID := c.Params("orgID")

				var organization models.Organization
				if err := initializers.DB.Where("id = ?", orgID).First(&organization).Error; err != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}

				picName, err := utils.UploadImage(c, "coverPic", helpers.ChatClient, 720, 720)
				if err != nil {
					return err
				}

				chat = models.GroupChat{
					UserID:         parsedLoggedInUserID,
					Title:          reqBody.Title,
					Description:    reqBody.Description,
					OrganizationID: &organization.ID,
				}

				if picName != "" {
					chat.CoverPic = picName
				}

				result := initializers.DB.Create(&chat)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}

				// orgMemberID := c.GetRespHeader("orgMemberID")
				// parsedID, _ := uuid.Parse(orgMemberID)
				// go routines.MarkProjectHistory(project.ID, parsedID, 8, nil, nil, nil, nil, nil)
			}

			for _, chatUserID := range chatUserIDs {
				parsedChatUserID, err := uuid.Parse(chatUserID)
				if err != nil {
					return &fiber.Error{Code: 500, Message: "Invalid User ID."} //TODO errors config for all types of error messages
				}

				groupChatMembership := models.GroupChatMembership{
					UserID:      parsedChatUserID,
					GroupChatID: chat.ID,
					Role:        models.ChatMember,
				}

				if parsedChatUserID == parsedLoggedInUserID {
					groupChatMembership.Role = models.ChatAdmin
				}

				result := initializers.DB.Create(&groupChatMembership)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}
			}

			if err := initializers.DB.Preload("Memberships").Preload("Memberships.User").Find(&chat, "id = ? ", chat.ID).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}
		}

		return c.Status(201).JSON(fiber.Map{
			"status":  "success",
			"message": "Chat Created",
			"chat":    chat,
		})
	}
}

func AddGroupChatMembers(chatType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqBody struct {
			UserIDs []string `json:"userIDs"`
		}

		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}
		chatUserIDs := reqBody.UserIDs
		groupChatID := c.Params("chatID")

		var chat models.GroupChat
		err := initializers.DB.First(&chat, "id = ?", groupChatID).Error
		if err != nil {
			return &fiber.Error{Code: 400, Message: "No chat of this id found."}
		}

		if chatType == "Group" {
			userID := c.GetRespHeader("loggedInUserID")
			parsedLoggedInUserID, err := uuid.Parse(userID)
			if err != nil {
				return &fiber.Error{Code: 500, Message: "Error Parsing the LoggedIn User ID."}
			}

			var invitations []models.Invitation

			for _, chatUserID := range chatUserIDs {
				parsedUserID, err := uuid.Parse(chatUserID)
				if err != nil {
					return &fiber.Error{Code: 500, Message: "Error Parsing the User ID."}
				}
				if parsedUserID == parsedLoggedInUserID {
					continue
				}

				var existingMembership models.GroupChatMembership
				err = initializers.DB.First(&existingMembership, "group_chat_id = ? AND user_id = ?", groupChatID, parsedUserID).Error
				if err == nil {
					continue
				}

				var existingInvitation models.Invitation
				err = initializers.DB.First(&existingInvitation, "group_chat_id = ? AND user_id = ? AND status=0", groupChatID, parsedUserID).Error
				if err == nil {
					continue
				}

				invitation := models.Invitation{
					UserID:      parsedUserID,
					GroupChatID: &chat.ID,
				}

				result := initializers.DB.Create(&invitation)
				if result.Error != nil {
					return &fiber.Error{Code: 500, Message: "Internal Server Error while creating invitations"}
				}

				invitations = append(invitations, invitation)
			}

			return c.Status(200).JSON(fiber.Map{
				"status":      "success",
				"message":     "Invitations Sent",
				"invitations": invitations,
			})
		} else {
			var memberships []models.GroupChatMembership

			for _, chatUserID := range chatUserIDs {
				parsedUserID, err := uuid.Parse(chatUserID)
				if err != nil {
					return &fiber.Error{Code: 500, Message: "Error Parsing the User ID."}
				}

				var existingMembership models.GroupChatMembership
				err = initializers.DB.First(&existingMembership, "group_chat_id = ? AND user_id = ?", groupChatID, parsedUserID).Error
				if err == nil {
					continue
				}

				var projectMembership models.Membership
				err = initializers.DB.First(&projectMembership, "project_id = ? AND user_id = ?", chat.ProjectID, parsedUserID).Error
				if err != nil {
					continue
				}

				groupChatMembership := models.GroupChatMembership{
					UserID:      parsedUserID,
					GroupChatID: chat.ID,
					Role:        models.ChatMember,
				}

				result := initializers.DB.Create(&groupChatMembership)
				if result.Error != nil {
					return &fiber.Error{Code: 500, Message: "Internal Server Error while creating memberships"}
				}

				memberships = append(memberships, groupChatMembership)
			}

			return c.Status(200).JSON(fiber.Map{
				"status":      "success",
				"message":     "Invitations Sent",
				"memberships": memberships,
			})
		}
	}
}

func RemoveGroupChatMember(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	var reqBody struct {
		UserID string `json:"userID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var chatMembership models.GroupChatMembership
	if err := initializers.DB.First(&chatMembership, "user_id = ? AND group_chat_id=?", reqBody.UserID, parsedChatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&chatMembership).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully",
	})
}

func EditGroupChat(c *fiber.Ctx) error {
	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AdminOnly   bool   `json:"adminOnly"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	groupChatID := c.Params("chatID")

	var groupChat models.GroupChat
	err := initializers.DB.First(&groupChat, "id = ?", groupChatID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No chat of this id found."}
	}

	if reqBody.Title != "" {
		groupChat.Title = reqBody.Title
	}
	if reqBody.Description != "" {
		groupChat.Description = reqBody.Description
	}
	groupChat.AdminOnly = reqBody.AdminOnly

	// picName, err := utils.SaveFile(c, "coverPic", "chat", true, 720, 720)
	picName, err := utils.UploadImage(c, "coverPic", helpers.ChatClient, 720, 720)
	if err != nil {
		return err
	}

	oldGroupPic := groupChat.CoverPic

	if picName != "" {
		groupChat.CoverPic = picName
	}

	result := initializers.DB.Save(&groupChat)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	if picName != "" {
		go routines.DeleteFromBucket(helpers.ChatClient, oldGroupPic)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Updated",
		"chat":    groupChat,
	})
}

func EditGroupChatRole(c *fiber.Ctx) error {
	var reqBody struct {
		UserID string               `json:"userID"`
		Role   models.GroupChatRole `json:"role"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	groupChatID := c.Params("chatID")

	var userChatMembership models.GroupChatMembership
	err := initializers.DB.First(&userChatMembership, "group_chat_id = ? AND user_id = ?", groupChatID, reqBody.UserID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "User is not a member of this chat."}
	}

	userChatMembership.Role = reqBody.Role
	result := initializers.DB.Save(&userChatMembership)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Membership Updated",
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&chat).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully",
	})
}

func LeaveGroupChat(c *fiber.Ctx) error { //! when no admin left then make the first joined member admin
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var membership models.GroupChatMembership
	if err := initializers.DB.First(&membership, "group_chat_id = ? AND user_id = ?", chatID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Chat Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&membership).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var groupChat models.GroupChat
	if err := initializers.DB.Preload("Memberships").First(&groupChat, "id = ?", chatID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if len(groupChat.Memberships) == 0 {
		if err := initializers.DB.Delete(&groupChat).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Group Chat left successfully",
	})
}
