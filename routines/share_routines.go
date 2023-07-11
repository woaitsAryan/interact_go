package routines

import (
	"log"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

func IncrementPostShare(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		log.Println("No Post of this ID found.")
	} else {
		post.NoShares++
		result := initializers.DB.Save(post)
		if result.Error != nil {
			log.Println("Database Error while updating Post.")
		}
	}
}

func IncrementProjectShare(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		log.Println("No Project of this ID found.")
	} else {
		project.NoShares++
		result := initializers.DB.Save(project)
		if result.Error != nil {
			log.Println("Database Error while updating Project.")
		}
	}
}
