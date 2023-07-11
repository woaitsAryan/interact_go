package routines

import (
	"log"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementCountsAndSendNotification(loggedInUserID uuid.UUID, toFollowID uuid.UUID) {
	var toFollowUser models.User
	err := initializers.DB.First(&toFollowUser, "id=?", toFollowID).Error
	if err != nil {
		log.Println("No User with this ID exists.")
	} else {
		var loggedInUser models.User
		err := initializers.DB.First(&loggedInUser, "id=?", loggedInUserID).Error

		if err != nil {
			log.Println("Error Retrieving User.")
		} else {
			toFollowUser.NoFollowers++
			if err := initializers.DB.Save(&toFollowUser).Error; err != nil {
				log.Println("Database Error while incrementing number followers.")
			}

			loggedInUser.NoFollowing++
			if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
				log.Println("Database Error while incrementing number following.")
			}

			notification := models.Notification{
				NotificationType: 0,
				UserID:           toFollowUser.ID,
				SenderID:         loggedInUserID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				log.Println("Database Error while creating notification.")
			}
		}
	}
}

func DecrementCounts(loggedInUserID uuid.UUID, toUnFollowID uuid.UUID) {
	var toUnFollowUser models.User
	err := initializers.DB.First(&toUnFollowUser, "id=?", toUnFollowID).Error
	if err != nil {
		log.Println("No User with this ID exists.")
	} else {
		var loggedInUser models.User
		err := initializers.DB.First(&loggedInUser, "id=?", loggedInUserID).Error

		if err != nil {
			log.Println("Error Retrieving User.")
		} else {
			toUnFollowUser.NoFollowers--
			if err := initializers.DB.Save(&toUnFollowUser).Error; err != nil {
				log.Println("Database Error while decrementing number followers.")
			}

			loggedInUser.NoFollowing--
			if err := initializers.DB.Save(&loggedInUser).Error; err != nil {
				log.Println("Database Error while decrementing number following.")
			}
		}
	}
}
