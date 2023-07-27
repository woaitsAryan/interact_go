package routines

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateLastViewed(userID uuid.UUID, projectID uuid.UUID) {

	var projectView models.LastViewed
	if err := initializers.DB.Preload("User").Where("user_id = ? AND project_id=?", userID, projectID).First(&projectView).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			projectView.ProjectID = projectID
			projectView.UserID = userID
			projectView.Timestamp = time.Now()

			result := initializers.DB.Create(&projectView)
			if result.Error != nil {
				helpers.LogDatabaseError("Error whiling creating last viewed-UpdateLastViewed", err, "go_routine")
			}
		} else {
			helpers.LogDatabaseError("Error whiling fetching last viewed-UpdateLastViewed", err, "go_routine")
		}
	} else {
		projectView.Timestamp = time.Now()

		result := initializers.DB.Save(&projectView)
		if result.Error != nil {
			helpers.LogDatabaseError("Error whiling updating last viewed-UpdateLastViewed", err, "go_routine")
		}
	}

}
