package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Pratham-Mishra04/interact/config"
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
			return
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
			return
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

func GetApplicationScore(application *models.Application) {
	var user models.User
	initializers.DB.First(&user, "id=?", application.UserID)

	var opening models.Opening
	initializers.DB.First(&opening, "id=?", application.OpeningID)

	var value_tags []string
	if opening.ProjectID != nil {
		var project models.Project
		initializers.DB.First(&project, "id=?", opening.ProjectID)
		value_tags = project.Tags
	} else {
		var org models.Organization
		initializers.DB.Preload("User").First(&org, "id=?", opening.OrganizationID)
		value_tags = org.User.Tags
	}

	reqBody, _ := json.Marshal(map[string]any{
		"cover_letter":               application.Content,
		"profile_topics":             user.Tags,
		"is_resume_included":         false,
		"resume_topics":              []string{},
		"opening_description_topics": opening.Tags,
		"organization_values_topics": value_tags,
		"years_of_experience":        application.YOE,
	})

	response, err := http.Post(initializers.CONFIG.ML_URL+config.APPLICATION_SCORE, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Print(err)
		helpers.LogServerError("Error while fetching Application Score-GetApplicationScore", err, "go_routine")
		return
	}
	defer response.Body.Close()

	var responseBody struct {
		Score float32 `json:"score"`
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		fmt.Print(err)
		helpers.LogServerError("Error while fetching Application Score-GetApplicationScore", err, "go_routine")
		return
	}

	application.Score = responseBody.Score

	if err := initializers.DB.Save(&application).Error; err != nil {
		helpers.LogDatabaseError("Error while saving Application-GetApplicationScore", err, "go_routine")
	}
}
