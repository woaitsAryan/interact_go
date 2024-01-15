package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementPostShare(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		helpers.LogDatabaseError("No Post of this ID found-IncrementPostShare.", err, "go_routine")
	} else {
		post.NoShares++
		result := initializers.DB.Save(post)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-IncrementPostShare", result.Error, "go_routine")
		}
	}
}

func IncrementAnnouncementShare(announcementID uuid.UUID) {
	var announcement models.Announcement
	if err := initializers.DB.First(&announcement, "id=?", announcementID).Error; err != nil {
		helpers.LogDatabaseError("No Post of this ID found-IncrementAnnouncementShare.", err, "go_routine")
	} else {
		announcement.NoShares++
		result := initializers.DB.Save(announcement)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-IncrementAnnouncementShare", result.Error, "go_routine")
		}
	}
}

func IncrementProjectShare(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		helpers.LogDatabaseError("No Project of this ID found-IncrementProjectShare.", err, "go_routine")
	} else {
		project.NoShares++
		result := initializers.DB.Save(project)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Project-IncrementProjectShare", result.Error, "go_routine")
		}
	}
}

func IncrementEventShare(eventID uuid.UUID) {
	var event models.Event
	if err := initializers.DB.First(&event, "id=?", eventID).Error; err != nil {
		helpers.LogDatabaseError("No Event of this ID found-IncrementEventShare.", err, "go_routine")
	} else {
		event.NoShares++
		result := initializers.DB.Save(event)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Event-IncrementEventShare", result.Error, "go_routine")
		}
	}
}
