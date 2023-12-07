package messaging_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func ShareItem(shareType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		parsedUserID, _ := uuid.Parse(loggedInUserID)

		var reqBody struct {
			Content   string         `json:"content"`
			Chats     pq.StringArray `json:"chats"`
			PostID    string         `json:"postID"`
			ProjectID string         `json:"projectID"`
			EventID   string         `json:"eventID"`
			OpeningID string         `json:"openingID"`
			ProfileID string         `json:"profileID"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		chats := reqBody.Chats

		for _, chatID := range chats {
			message := models.Message{
				UserID:  parsedUserID,
				Content: reqBody.Content,
			}

			parsedChatID, err := uuid.Parse(chatID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID."}
			}

			var chat models.Chat
			if err := initializers.DB.Where("id=?", parsedChatID).First(&chat).Error; err != nil {
				continue
			}

			if parsedUserID == chat.AcceptingUserID && chat.BlockedByCreatingUser {
				continue
			}

			if parsedUserID == chat.CreatingUserID && chat.BlockedByAcceptingUser {
				continue
			}

			message.ChatID = chat.ID

			switch shareType {
			case "post":
				parsedPostID, err := uuid.Parse(reqBody.PostID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
				}
				message.PostID = &parsedPostID
				go routines.IncrementPostShare(parsedPostID)
			case "project":
				parsedProjectID, err := uuid.Parse(reqBody.ProjectID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
				}
				message.ProjectID = &parsedProjectID
				go routines.IncrementProjectShare(parsedProjectID)
			case "event":
				parsedEventID, err := uuid.Parse(reqBody.EventID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Event ID."}
				}
				message.EventID = &parsedEventID
				go routines.IncrementEventShare(parsedEventID)
			case "opening":
				parsedOpeningID, err := uuid.Parse(reqBody.OpeningID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Opening ID."}
				}
				message.OpeningID = &parsedOpeningID
			case "profile":
				parsedProfileID, err := uuid.Parse(reqBody.ProfileID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Profile ID."}
				}
				message.ProfileID = &parsedProfileID
			default:
				return &fiber.Error{Code: 400, Message: "Invalid Share Type."}
			}

			result := initializers.DB.Create(&message)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Shared",
		})

	}
}
