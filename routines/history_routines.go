package routines

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateLastViewedProject(userID uuid.UUID, projectID uuid.UUID) {
	var projectView models.LastViewedProjects
	if err := initializers.DB.Preload("User").Where("user_id = ? AND project_id=?", userID, projectID).First(&projectView).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			projectView.ProjectID = projectID
			projectView.UserID = userID
			projectView.Timestamp = time.Now()

			result := initializers.DB.Create(&projectView)
			if result.Error != nil {
				helpers.LogDatabaseError("Error whiling creating last viewed-UpdateLastViewedProject", result.Error, "go_routine")
			}
		} else {
			helpers.LogDatabaseError("Error whiling fetching last viewed-UpdateLastViewedProject", err, "go_routine")
		}
	} else {
		projectView.Timestamp = time.Now()

		result := initializers.DB.Save(&projectView)
		if result.Error != nil {
			helpers.LogDatabaseError("Error whiling updating last viewed-UpdateLastViewedProject", result.Error, "go_routine")
		}
	}
}

func UpdateLastViewedOpening(userID uuid.UUID, openingID uuid.UUID) {
	var openingView models.LastViewedOpenings
	if err := initializers.DB.Preload("User").Where("user_id = ? AND project_id=?", userID, openingID).First(&openingView).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			openingView.OpeningID = openingID
			openingView.UserID = userID
			openingView.Timestamp = time.Now()

			result := initializers.DB.Create(&openingView)
			if result.Error != nil {
				helpers.LogDatabaseError("Error whiling creating last viewed-UpdateLastViewedOpening", result.Error, "go_routine")
			}
		} else {
			helpers.LogDatabaseError("Error whiling fetching last viewed-UpdateLastViewedOpening", err, "go_routine")
		}
	} else {
		openingView.Timestamp = time.Now()

		result := initializers.DB.Save(&openingView)
		if result.Error != nil {
			helpers.LogDatabaseError("Error whiling updating last viewed-UpdateLastViewedOpening", result.Error, "go_routine")
		}
	}
}

func MarkProjectHistory(
	projectID uuid.UUID,
	senderID uuid.UUID,
	historyType int,
	userID *uuid.UUID,
	openingID *uuid.UUID,
	applicationID *uuid.UUID,
	invitationID *uuid.UUID,
	taskID *uuid.UUID) {

	history := models.ProjectHistory{
		ProjectID:   projectID,
		SenderID:    senderID,
		HistoryType: historyType,
	}

	if userID != nil {
		history.UserID = userID
	}
	if openingID != nil {
		history.OpeningID = openingID
	}
	if applicationID != nil {
		history.ApplicationID = applicationID
	}
	if invitationID != nil {
		history.InvitationID = invitationID
	}
	if taskID != nil {
		history.TaskID = taskID
	}

	if err := initializers.DB.Create(&history).Error; err != nil {
		helpers.LogDatabaseError("Error while creating Project History-MarkProjectHistory", err, "go_routine")
	}
}

func MarkOrganizationHistory(
	orgID uuid.UUID,
	userID uuid.UUID,
	historyType int,
	postID *uuid.UUID,
	projectID *uuid.UUID,
	eventID *uuid.UUID,
	taskID *uuid.UUID,
	invitationID *uuid.UUID) {
	
	organizationHistory := models.OrganizationHistory{
		HistoryType: historyType,
		OrganizationID: orgID,
		UserID: userID,
	}

	if postID != nil {
		organizationHistory.PostID = postID
	}
	if projectID != nil {
		organizationHistory.ProjectID = projectID
	}
	if eventID != nil {
		organizationHistory.EventID = eventID
	}
	if taskID != nil {
		organizationHistory.TaskID = taskID
	}
	if invitationID != nil {
		organizationHistory.InvitationID = invitationID
	}
	if err := initializers.DB.Create(&organizationHistory).Error; err != nil {
		helpers.LogDatabaseError("Error while creating Organization History-MarkOrganizationHistory", err, "go_routine")
	}
}