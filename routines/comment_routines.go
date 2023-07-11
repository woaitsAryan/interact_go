package routines

import (
	"log"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementPostCommentsAndSendNotification(postID uuid.UUID, loggedInUserID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		log.Println("No Post of this ID found.")
	} else {
		post.NoComments++
		result := initializers.DB.Save(&post)

		if result.Error != nil {
			log.Println("Internal Server Error while saving the post.")
		}

		if loggedInUserID != post.UserID {
			notification := models.Notification{
				SenderID:         loggedInUserID,
				NotificationType: 2,
				UserID:           post.UserID,
				PostID:           &post.ID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				log.Println("Database Error while creating notification.")
			}
		}
	}
}

func IncrementProjectCommentsAndSendNotification(projectID uuid.UUID, loggedInUserID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		log.Println("No Project of this ID found.")
	} else {
		project.NoComments++
		result := initializers.DB.Save(&project)

		if result.Error != nil {
			log.Println("Internal Server Error while saving the project.")
		}

		if loggedInUserID != project.UserID {
			notification := models.Notification{
				SenderID:         loggedInUserID,
				NotificationType: 4,
				UserID:           project.UserID,
				ProjectID:        &project.ID,
			}

			if err := initializers.DB.Create(&notification).Error; err != nil {
				log.Println("Database Error while creating notification.")
			}
		}
	}
}

func DecrementPostComments(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		log.Println("No Post of this ID found.")
	} else {
		post.NoComments--
		result := initializers.DB.Save(&post)

		if result.Error != nil {
			log.Println("Internal Server Error while saving the post.")
		}
	}
}

func DecrementProjectComments(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		log.Println("No Project of this ID found.")
	} else {
		project.NoComments--
		result := initializers.DB.Save(&project)

		if result.Error != nil {
			log.Println("Internal Server Error while saving the project.")
		}
	}
}
