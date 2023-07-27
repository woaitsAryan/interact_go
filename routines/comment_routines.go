package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementPostCommentsAndSendNotification(postID uuid.UUID, loggedInUserID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		helpers.LogDatabaseError("No Post of this ID found-IncrementPostCommentsAndSendNotification.", err, "go_routine")
	} else {
		post.NoComments++
		result := initializers.DB.Save(&post)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-IncrementPostCommentsAndSendNotification", err, "go_routine")
		}

		if loggedInUserID != post.UserID {
			notification := models.Notification{
				SenderID:         loggedInUserID,
				NotificationType: 2,
				UserID:           post.UserID,
				PostID:           &post.ID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				helpers.LogDatabaseError("Error while creating Notification-IncrementPostCommentsAndSendNotification", err, "go_routine")
			}
		}
	}
}

func IncrementProjectCommentsAndSendNotification(projectID uuid.UUID, loggedInUserID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		helpers.LogDatabaseError("No Project of this ID found-IncrementProjectCommentsAndSendNotification.", err, "go_routine")
	} else {
		project.NoComments++
		result := initializers.DB.Save(&project)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Project-IncrementProjectCommentsAndSendNotification", err, "go_routine")
		}

		if loggedInUserID != project.UserID {
			notification := models.Notification{
				SenderID:         loggedInUserID,
				NotificationType: 4,
				UserID:           project.UserID,
				ProjectID:        &project.ID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				helpers.LogDatabaseError("Error while creating Notification-IncrementProjectCommentsAndSendNotification", err, "go_routine")
			}
		}
	}
}

func DecrementPostComments(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		helpers.LogDatabaseError("No Post of this ID found-DecrementPostComments.", err, "go_routine")
	} else {
		post.NoComments--
		result := initializers.DB.Save(&post)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-DecrementPostComments", err, "go_routine")
		}
	}
}

func DecrementProjectComments(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		helpers.LogDatabaseError("No Project of this ID found-DecrementProjectComments.", err, "go_routine")
	} else {
		project.NoComments--
		result := initializers.DB.Save(&project)

		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Project-DecrementProjectComments", err, "go_routine")
		}
	}
}
