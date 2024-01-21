package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementOpeningApplicationsAndSendNotification(openingID uuid.UUID, applicationID uuid.UUID, userID uuid.UUID) {
	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", openingID).Error; err != nil {
		helpers.LogDatabaseError("No Opening of this ID found-IncrementOpeningApplicationsAndSendNotification.", err, "go_routine")
	} else {
		opening.NoOfApplications++
		result := initializers.DB.Save(&opening)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Opening-IncrementOpeningApplicationsAndSendNotification", err, "go_routine")
		}

		notification := models.Notification{
			NotificationType: 5,
			UserID:           opening.UserID,
			SenderID:         userID,
			OpeningID:        &opening.ID,
			ApplicationID:    &applicationID,
			ProjectID:        opening.ProjectID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			helpers.LogDatabaseError("Error while creating Notification-IncrementOpeningApplicationsAndSendNotification", err, "go_routine")
		}
	}

}

func DecrementOpeningApplications(openingID uuid.UUID) {
	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", openingID).Error; err != nil {
		helpers.LogDatabaseError("No Opening of this ID found-DecrementOpeningApplications.", err, "go_routine")
	} else {
		opening.NoOfApplications--
		result := initializers.DB.Save(&opening)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Opening-DecrementOpeningApplications", err, "go_routine")
		}
	}

}

func ProjectMembershipSendNotification(application *models.Application) {

	notification := models.Notification{
		NotificationType: 6,
		UserID:           application.UserID,
		OpeningID:        &application.OpeningID,
	}

	if err := initializers.DB.Create(&notification).Error; err != nil {
		helpers.LogDatabaseError("Error while creating Notification-CreateMembershipAndSendNotification", err, "go_routine")
	}
}

func IncrementOrgOpeningApplicationsAndSendNotification(openingID uuid.UUID, applicationID uuid.UUID, userID uuid.UUID) {
	var opening models.Opening
	if err := initializers.DB.First(&opening, "id=?", openingID).Error; err != nil {
		helpers.LogDatabaseError("No Opening of this ID found-IncrementOrgOpeningApplicationsAndSendNotification.", err, "go_routine")
	} else {
		opening.NoOfApplications++
		result := initializers.DB.Save(&opening)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Opening-IncrementOrgOpeningApplicationsAndSendNotification", err, "go_routine")
		}

		notification := models.Notification{
			NotificationType: 20,
			UserID:           opening.UserID,
			SenderID:         userID,
			OpeningID:        &opening.ID,
			ApplicationID:    &applicationID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			helpers.LogDatabaseError("Error while creating Notification-IncrementOrgOpeningApplicationsAndSendNotification", err, "go_routine")
		}
	}
}

func OrgMembershipSendNotification(application *models.Application) {
	
	notification := models.Notification{
		NotificationType: 21,
		UserID:           application.UserID,
		OpeningID:        &application.OpeningID,
	}

	if err := initializers.DB.Create(&notification).Error; err != nil {
		helpers.LogDatabaseError("Error while creating Notification-CreateOrgMembershipAndSendNotification", err, "go_routine")
	}
}
