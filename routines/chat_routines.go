package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func UpdateChatLastRead(chatID uuid.UUID, messages []models.Message, loggedInUserID uuid.UUID) {
	var chat models.Chat
	if err := initializers.DB.Preload("LastReadMessageByCreatingUser").
		Preload("LastReadMessageByAcceptingUser").
		First(&chat, "id=?", chatID).Error; err != nil {
		helpers.LogDatabaseError("Error while fetching Chat-UpdateChatLastRead", err, "go_routine")
	}

	if chat.AcceptingUserID == loggedInUserID {
		for _, msg := range messages {
			if msg.UserID.String() == chat.CreatingUserID.String() {
				if msg.CreatedAt.After(chat.LastReadMessageByAcceptingUser.CreatedAt) {
					chat.LastReadMessageByAcceptingUserID = msg.ID

					result := initializers.DB.Save(&chat)
					if result.Error != nil {
						helpers.LogDatabaseError("Error while updating Chat-UpdateChatLastRead", result.Error, "go_routine")
					}
				}
				break
			}
		}
	} else if chat.CreatingUserID == loggedInUserID {
		for _, msg := range messages {
			if msg.UserID.String() == chat.AcceptingUserID.String() {
				if msg.CreatedAt.After(chat.LastReadMessageByCreatingUser.CreatedAt) {
					chat.LastReadMessageByCreatingUserID = msg.ID

					result := initializers.DB.Save(&chat)
					if result.Error != nil {
						helpers.LogDatabaseError("Error while updating Chat-UpdateChatLastRead", result.Error, "go_routine")
					}
				}
				break
			}
		}
	}
}
