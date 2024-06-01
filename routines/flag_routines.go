package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/utils"
)

func CheckFlagComment(comment *models.Comment) {
	flag, err := utils.MLFlagReq(comment.Content)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, "CheckFlagComment")
	} else {
		comment.IsFlagged = flag

		if err := initializers.DB.Save(&comment).Error; err != nil {
			helpers.LogDatabaseError("Error while saving Comment-CheckFlagComment", err, "go_routine")
		}

		//TODO send toxicity/flagged mail
	}
}

func CheckFlagPost(post *models.Post) {
	flag, err := utils.MLFlagReq(post.Content)
	if err != nil {
		helpers.LogServerError("Error Fetching from ML API", err, "CheckFlagPost")
	} else {
		post.IsFlagged = flag

		if err := initializers.DB.Save(&post).Error; err != nil {
			helpers.LogDatabaseError("Error while saving Post-CheckFlagPost", err, "go_routine")
		}

		//TODO send toxicity/flagged mail
	}
}
