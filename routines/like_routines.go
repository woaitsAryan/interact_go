package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func IncrementPostLikesAndSendNotification(postID uuid.UUID, loggedInUserID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post of this ID found-IncrementPostLikesAndSendNotification.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post-IncrementPostLikesAndSendNotification", err, "go_routine")
		}
	} else {
		post.NoLikes++

		result := initializers.DB.Save(&post)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-IncrementPostLikesAndSendNotification", result.Error, "go_routine")
		}
	}

	if loggedInUserID != post.UserID {

		notification := models.Notification{
			NotificationType: 1,
			UserID:           post.UserID,
			SenderID:         loggedInUserID,
			PostID:           &post.ID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			helpers.LogDatabaseError("Error while creating Notification-IncrementPostLikesAndSendNotification", err, "go_routine")
		}
	}
}

func DecrementPostLikes(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post of this ID found-DecrementPostLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post-DecrementPostLikes", err, "go_routine")
		}
	} else {
		post.NoLikes--

		result := initializers.DB.Save(&post)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-DecrementPostLikes", result.Error, "go_routine")
		}
	}
}

func IncrementProjectLikesAndSendNotification(projectID uuid.UUID, loggedInUserID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Project of this ID found-IncrementProjectLikesAndSendNotification.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Project-IncrementProjectLikesAndSendNotification", err, "go_routine")
		}
	} else {
		project.NoLikes++

		result := initializers.DB.Save(&project)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Project-IncrementProjectLikesAndSendNotification", result.Error, "go_routine")
		}
	}

	if loggedInUserID != project.UserID {

		notification := models.Notification{
			NotificationType: 3,
			UserID:           project.UserID,
			SenderID:         loggedInUserID,
			ProjectID:        &project.ID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			helpers.LogDatabaseError("Error while creating Notification-IncrementProjectLikesAndSendNotification", err, "go_routine")
		}
	}
}

func DecrementProjectLikes(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Project of this ID found-DecrementProjectLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Project-DecrementProjectLikes", err, "go_routine")
		}
	} else {
		project.NoLikes--

		result := initializers.DB.Save(&project)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Project-DecrementProjectLikes", result.Error, "go_routine")
		}
	}
}

func IncrementCommentLikes(commentID uuid.UUID, loggedInUserID uuid.UUID) {
	var comment models.Comment
	if err := initializers.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post Comment of this ID found-IncrementPostCommentLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post Comment-IncrementPostCommentLikes", err, "go_routine")
		}
	} else {
		comment.NoLikes++

		result := initializers.DB.Save(&comment)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post Comment-IncrementPostCommentLikes", result.Error, "go_routine")
		}
	}
}

func DecrementCommentLikes(commentID uuid.UUID) {
	var comment models.Comment
	if err := initializers.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post Comment of this ID found-DecrementPostCommentLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post Comment-DecrementPostCommentLikes", err, "go_routine")
		}
	} else {
		comment.NoLikes--

		result := initializers.DB.Save(&comment)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post Comment-DecrementPostCommentLikes", result.Error, "go_routine")
		}
	}
}

func IncrementEventLikesAndSendNotification(eventID uuid.UUID, loggedInUserID uuid.UUID) {
	var event models.Event
	if err := initializers.DB.Preload("Organization").First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Event of this ID found-IncrementEventLikesAndSendNotification.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Event-IncrementEventLikesAndSendNotification", err, "go_routine")
		}
	} else {
		event.NoLikes++

		result := initializers.DB.Save(&event)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Event-IncrementEventLikesAndSendNotification", result.Error, "go_routine")
		}
	}

	if loggedInUserID != event.Organization.UserID {
		notification := models.Notification{
			NotificationType: 12,
			UserID:           event.Organization.UserID,
			SenderID:         loggedInUserID,
			EventID:          &event.ID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			helpers.LogDatabaseError("Error while creating Notification-IncrementEventLikesAndSendNotification", err, "go_routine")
		}
	}
}

func DecrementEventLikes(eventID uuid.UUID) {
	var event models.Event
	if err := initializers.DB.First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Event of this ID found-DecrementEventLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Event-DecrementEventLikes", err, "go_routine")
		}
	} else {
		event.NoLikes--

		result := initializers.DB.Save(&event)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Event-DecrementEventLikes", result.Error, "go_routine")
		}
	}
}

func IncrementReviewLikes(eventID uuid.UUID, loggedInUserID uuid.UUID) {
	var review models.Review
	if err := initializers.DB.Preload("Organization").First(&review, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Review of this ID found-IncrementReviewLikesAndSendNotification.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Review-IncrementReviewLikesAndSendNotification", err, "go_routine")
		}
	} else {
		review.NumberOfUpVotes++

		result := initializers.DB.Save(&review)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Review-IncrementReviewLikesAndSendNotification", result.Error, "go_routine")
		}
	}

	// if loggedInUserID != review.Organization.UserID && loggedInUserID != review.UserID {
	// 	notification := models.Notification{
	// 		NotificationType: 12,
	// 		UserID:           review.Organization.UserID,
	// 		SenderID:         loggedInUserID,
	// 		EventID:          &review.ID,
	// 	}

	// 	if err := initializers.DB.Create(&notification).Error; err != nil {
	// 		helpers.LogDatabaseError("Error while creating Notification-IncrementReviewLikesAndSendNotification", err, "go_routine")
	// 	}
	// }
}

func DecrementReviewLikes(reviewID uuid.UUID) {
	var review models.Review
	if err := initializers.DB.First(&review, "id = ?", reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Review of this ID found-DecrementReviewLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Review-DecrementReviewLikes", err, "go_routine")
		}
	} else {
		review.NumberOfUpVotes--

		result := initializers.DB.Save(&review)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Review-DecrementReviewLikes", result.Error, "go_routine")
		}
	}
}
