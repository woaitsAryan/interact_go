package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func SendWelcomeNotification(userID uuid.UUID) {
	notification := models.Notification{
		NotificationType: -1,
		UserID:           userID,
		SenderID:         userID,
	}
	result := initializers.DB.Create(&notification)
	if result.Error != nil {
		helpers.LogDatabaseError("Error whiling creating notification-AddWelcomeNotification", result.Error, "go_routine")
	}
}

func SendChatNotification(creatorID uuid.UUID, acceptorID uuid.UUID) {
	notification := models.Notification{
		NotificationType: 9,
		UserID:           acceptorID,
		SenderID:         creatorID,
	}
	result := initializers.DB.Create(&notification)
	if result.Error != nil {
		helpers.LogDatabaseError("Error whiling creating notification-AddWelcomeNotification", result.Error, "go_routine")
	}
}
