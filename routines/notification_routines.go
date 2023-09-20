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
func SendInvitationAcceptedNotification(creatorID uuid.UUID, acceptorID uuid.UUID) {
	notification := models.Notification{
		NotificationType: 10,
		UserID:           creatorID,
		SenderID:         acceptorID,
	}
	result := initializers.DB.Create(&notification)
	if result.Error != nil {
		helpers.LogDatabaseError("Error whiling creating notification-AddWelcomeNotification", result.Error, "go_routine")
	}
}

func MarkReadNotifications(UnreadNotifications []uuid.UUID) {
	for _, unreadNotificationID := range UnreadNotifications {
		var notification models.Notification
		if err := initializers.DB.
			Where("id=?", unreadNotificationID).
			First(&notification).
			Error; err != nil {
			helpers.LogDatabaseError("Error whiling creating notification-AddWelcomeNotification", err, "go_routine")
		}
		notification.Read = true
		result := initializers.DB.Save(&notification)
		if result.Error != nil {
			helpers.LogDatabaseError("Error whiling creating notification-AddWelcomeNotification", result.Error, "go_routine")
		}
	}
}
