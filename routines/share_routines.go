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
			helpers.LogDatabaseError("Error while updating Post-IncrementPostShare", err, "go_routine")
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
			helpers.LogDatabaseError("Error while updating Project-IncrementProjectShare", err, "go_routine")
		}
	}
}
