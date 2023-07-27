package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementCountsAndSendNotification(loggedInUserID uuid.UUID, toFollowID uuid.UUID) {
	var toFollowUser models.User
	err := initializers.DB.First(&toFollowUser, "id=?", toFollowID).Error
	if err != nil {
		helpers.LogDatabaseError("No User of this ID found-IncrementCountsAndSendNotification.", err, "go_routine")
	} else {
		var loggedInUser models.User
		err := initializers.DB.First(&loggedInUser, "id=?", loggedInUserID).Error

		if err != nil {
			helpers.LogDatabaseError("No User of this LoggedIn ID found-IncrementCountsAndSendNotification.", err, "go_routine")
		} else {
			toFollowUser.NoFollowers++
			if err := initializers.DB.Save(&toFollowUser).Error; err != nil {
				helpers.LogDatabaseError("Error while incrementing number followers-IncrementCountsAndSendNotification", err, "go_routine")
			}

			loggedInUser.NoFollowing++
			if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
				helpers.LogDatabaseError("Error while incrementing number following-IncrementCountsAndSendNotification", err, "go_routine")
			}

			notification := models.Notification{
				NotificationType: 0,
				UserID:           toFollowUser.ID,
				SenderID:         loggedInUserID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				helpers.LogDatabaseError("Error while creating Notification-IncrementCountsAndSendNotification", err, "go_routine")
			}
		}
	}
}

func DecrementCounts(loggedInUserID uuid.UUID, toUnFollowID uuid.UUID) {
	var toUnFollowUser models.User
	err := initializers.DB.First(&toUnFollowUser, "id=?", toUnFollowID).Error
	if err != nil {
		helpers.LogDatabaseError("No User of this ID found-DecrementCounts.", err, "go_routine")
	} else {
		var loggedInUser models.User
		err := initializers.DB.First(&loggedInUser, "id=?", loggedInUserID).Error

		if err != nil {
			helpers.LogDatabaseError("No User of this LoggedIn ID found-DecrementCounts.", err, "go_routine")
		} else {
			toUnFollowUser.NoFollowers--
			if err := initializers.DB.Save(&toUnFollowUser).Error; err != nil {
				helpers.LogDatabaseError("Error while decrementing number followers-DecrementCounts", err, "go_routine")
			}

			loggedInUser.NoFollowing--
			if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
				helpers.LogDatabaseError("Error while decrementing number following-DecrementCounts", err, "go_routine")
			}
		}
	}
}
