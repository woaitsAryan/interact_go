package routines

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func UpdateProfileViews(user *models.User) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var profileView models.ProfileView
	initializers.DB.Where("user_id = ? AND date = ?", user.ID, today).First(&profileView)

	if profileView.ID == uuid.Nil {
		profileView = models.ProfileView{
			UserID: user.ID,
			Date:   today,
			Count:  1,
		}
		initializers.DB.Create(&profileView)
	} else {
		profileView.Count++
		initializers.DB.Save(&profileView)
	}

	user.TotalNoViews++
	result := initializers.DB.Save(user)
	if result.Error != nil {
		helpers.LogDatabaseError("Error while updating User-UpdateProfileViews", result.Error, "go_routine")
	}
}

func UpdateProjectViews(project *models.Project) { //TODO Creator and Member Check
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var projectView models.ProjectView
	initializers.DB.Where("project_id = ? AND date = ?", project.ID, today).First(&projectView)

	if projectView.ID == uuid.Nil {
		projectView = models.ProjectView{
			ProjectID: project.ID,
			Date:      today,
			Count:     1,
		}
		initializers.DB.Create(&projectView)
	} else {
		projectView.Count++
		initializers.DB.Save(&projectView)
	}

	project.TotalNoViews++
	result := initializers.DB.Save(project)
	if result.Error != nil {
		helpers.LogDatabaseError("Error while updating Project-UpdateProjectViews", result.Error, "go_routine")
	}

}
