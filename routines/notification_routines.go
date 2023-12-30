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

func SendTaskNotification(userID uuid.UUID, senderID uuid.UUID, projectID uuid.UUID) {
	notification := models.Notification{
		NotificationType: 11,
		UserID:           userID,
		SenderID:         senderID,
		ProjectID:        &projectID,
	}
	result := initializers.DB.Create(&notification)
	if result.Error != nil {
		helpers.LogDatabaseError("Error whiling creating notification-SendTaskNotification", result.Error, "go_routine")
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

func sendImpressionNotification(userID uuid.UUID,
	senderID uuid.UUID,
	postID *uuid.UUID,
	projectID *uuid.UUID,
	eventID *uuid.UUID,
	impressionCount int) {

	if !(impressionCount == 50 || impressionCount == 200 || impressionCount == 500 ||
		impressionCount == 1000 || impressionCount%1000 == 0) {
		return
	}

	notification := models.Notification{
		UserID:          userID,
		SenderID:        senderID,
		ImpressionCount: impressionCount,
	}
	if postID != nil {
		notification.PostID = postID
		notification.NotificationType = 14
	}
	if projectID != nil {
		notification.ProjectID = projectID
		notification.NotificationType = 15
	}
	if eventID != nil {
		notification.EventID = eventID
		notification.NotificationType = 16
	}

	result := initializers.DB.Create(&notification)
	if result.Error != nil {
		helpers.LogDatabaseError("Error whiling creating notification-SendImpressionNotification", result.Error, "go_routine")
	}
}
