package routines

import (
	"log"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementOpeningApplicationsAndSendNotification(openingID uuid.UUID, applicationID uuid.UUID, userID uuid.UUID) {
	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", openingID).Error; err != nil {
		log.Println("No Opening of this ID found.")
	} else {
		opening.NoOfApplications++
		result := initializers.DB.Save(&opening)

		if result.Error != nil {
			log.Println("Internal Server Error while saving the opening.")
		}

		notification := models.Notification{
			NotificationType: 5,
			UserID:           opening.UserID,
			SenderID:         userID,
			OpeningID:        &opening.ID,
			ApplicationID:    &applicationID,
			ProjectID:        &opening.ProjectID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			log.Println("Database Error while creating notification.")
		}
	}

}

func DecrementOpeningApplications(openingID uuid.UUID) {
	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", openingID).Error; err != nil {
		log.Println("No Opening of this ID found.")
	} else {
		opening.NoOfApplications--
		result := initializers.DB.Save(&opening)

		if result.Error != nil {
			log.Println("Internal Server Error while saving the opening.")
		}
	}

}

func CreateMembershipAndSendNotification(application *models.Application) {
	membership := models.Membership{
		ProjectID: application.Opening.ProjectID,
		UserID:    application.UserID,
		Role:      "",
		Title:     application.Opening.Title,
	}

	result := initializers.DB.Create(&membership)

	if result.Error != nil {
		log.Println("Internal Server Error while creating membership.")
	}

	notification := models.Notification{
		NotificationType: 6,
		UserID:           application.UserID,
		OpeningID:        &application.OpeningID,
	}

	if err := initializers.DB.Create(&notification).Error; err != nil {
		log.Println("Internal Server Error while creating notification.")
	}
}
