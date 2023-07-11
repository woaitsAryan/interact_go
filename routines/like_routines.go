package routines

import (
	"log"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func IncrementPostLikesAndSendNotification(postID uuid.UUID, loggedInUserID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Post of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		post.NoLikes++

		result := initializers.DB.Save(&post)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the post.")
		}
	}

	if loggedInUserID != post.UserID {

		notification := models.Notification{
			NotificationType: 3,
			UserID:           post.UserID,
			SenderID:         loggedInUserID,
			PostID:           &post.ID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			log.Println("Internal Server Error while creating the notification.")
		}
	}
}

func DecrementPostLikes(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Post of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		post.NoLikes--

		result := initializers.DB.Save(&post)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the post.")
		}
	}
}

func IncrementProjectLikesAndSendNotification(projectID uuid.UUID, loggedInUserID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Project of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		project.NoLikes++

		result := initializers.DB.Save(&project)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the project.")
		}
	}

	if loggedInUserID != project.UserID {

		notification := models.Notification{
			NotificationType: 1,
			UserID:           project.UserID,
			SenderID:         loggedInUserID,
			ProjectID:        &project.ID,
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			log.Println("Internal Server Error while creating the notification.")
		}
	}
}

func DecrementProjectLikes(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Project of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		project.NoLikes--

		result := initializers.DB.Save(&project)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the project.")
		}
	}
}

func IncrementPostCommentLikes(commentID uuid.UUID, loggedInUserID uuid.UUID) {
	var comment models.PostComment
	if err := initializers.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Comment of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		comment.NoLikes++

		result := initializers.DB.Save(&comment)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the comment.")
		}
	}
}

func DecrementPostCommentLikes(commentID uuid.UUID) {
	var comment models.PostComment
	if err := initializers.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Comment of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		comment.NoLikes--

		result := initializers.DB.Save(&comment)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the comment.")
		}
	}
}

func IncrementProjectCommentLikes(commentID uuid.UUID, loggedInUserID uuid.UUID) {
	var comment models.ProjectComment
	if err := initializers.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Comment of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		comment.NoLikes++

		result := initializers.DB.Save(&comment)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the comment.")
		}
	}
}

func DecrementProjectCommentLikes(commentID uuid.UUID) {
	var comment models.ProjectComment
	if err := initializers.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No Comment of this ID found.")
		} else {
			log.Println("Database Error.")
		}
	} else {
		comment.NoLikes--

		result := initializers.DB.Save(&comment)
		if result.Error != nil {
			log.Println("Internal Server Error while saving the comment.")
		}
	}
}
