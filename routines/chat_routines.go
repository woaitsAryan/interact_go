package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func UpdateChatLastRead(chatID uuid.UUID, loggedInUserID uuid.UUID) {
	var chat models.Chat
	if err := initializers.DB.First(&chat, "id=?", chatID).Error; err != nil {
		helpers.LogDatabaseError("Error while fetching Chat-UpdateChatLastRead", err, "go_routine")
	}

	// if chat.AcceptingUserID == loggedInUserID { //! Will have to update this for paginated db
	// 	for _, msg := range messages { // set a primary index on createdAt and fetch the first one
	// 		if msg.UserID.String() == chat.CreatingUserID.String() {
	// 			chat.LastReadMessageByAcceptingUserID = msg.ID
	// 			break
	// 		}
	// 	}
	// } else if chat.CreatingUserID == loggedInUserID { //! Will have to update this for paginated db
	// 	for _, msg := range messages {
	// 		if msg.UserID.String() == chat.AcceptingUserID.String() {
	// 			chat.LastReadMessageByCreatingUserID = msg.ID
	// 			break
	// 		}
	// 	}
	// }

	// result := initializers.DB.Save(&chat)
	// if result.Error != nil {
	// 	helpers.LogDatabaseError("Error while updating Chat-UpdateChatLastRead", result.Error, "go_routine")
	// }
}
