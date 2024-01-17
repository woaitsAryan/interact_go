package routines

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func IncrementOrgMember(orgID uuid.UUID) {
	var org models.Organization
	if err := initializers.DB.First(&org, "id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Org of this ID found-IncrementMemberCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Org-IncrementMemberCount", err, "go_routine")
		}
	} else {
		org.NumberOfMembers++

		result := initializers.DB.Save(&org)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Org-IncrementMemberCount", result.Error, "go_routine")
		}
	}
}

func IncrementOrgProject(orgID uuid.UUID) {
	var org models.Organization
	if err := initializers.DB.First(&org, "id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Org of this ID found-IncrementProjectCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Org-IncrementProjectCount", err, "go_routine")
		}
	} else {
		org.NumberOfProjects++
		if err := initializers.DB.Save(&org).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Org-IncrementProjecCount", err, "go_routine")
		}
	}
}

func IncrementOrgEvent(orgID uuid.UUID) {
	var org models.Organization
	if err := initializers.DB.First(&org, "id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Org of this ID found-IncrementEventCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Org-IncrementEventCount", err, "go_routine")
		}
	} else {
		org.NumberOfEvents++
		if err := initializers.DB.Save(&org).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Org-IncrementEventCount", err, "go_routine")
		}
	}
}

func DecrementOrgMember(orgID uuid.UUID) {
	var org models.Organization
	if err := initializers.DB.First(&org, "id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Org of this ID found-DecrementMemberCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Org-DecrementMemberCount", err, "go_routine")
		}
	} else {
		org.NumberOfMembers--
		if err := initializers.DB.Save(&org).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Org-DecrementMemberCount", err, "go_routine")
		}
	}
}

func DecrementOrgProject(orgID uuid.UUID) {
	var org models.Organization
	if err := initializers.DB.First(&org, "id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Org of this ID found-DecrementProjectCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Org-DecrementProjectCount", err, "go_routine")
		}
	} else {
		org.NumberOfProjects--
		if err := initializers.DB.Save(&org).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Org-DecrementProjectCount", err, "go_routine")
		}
	}
}

func DecrementOrgEvent(orgID uuid.UUID) {
	var org models.Organization
	if err := initializers.DB.First(&org, "id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Org of this ID found-DecrementEventCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Org-DecrementEventCount", err, "go_routine")
		}
	} else {
		org.NumberOfEvents--
		if err := initializers.DB.Save(&org).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Org-DecrementEventCount", err, "go_routine")
		}
	}
}

func IncrementProjectMember(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Project of this ID found-IncrementMemberCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Project-IncrementMemberCount", err, "go_routine")
		}
	} else {
		project.NumberOfMembers++
		if err := initializers.DB.Save(&project).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Project-IncrementMemberCount", err, "go_routine")
		}
		setUserCollaborativeProject(project.UserID)
	}
}

func DecrementProjectMember(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Project of this ID found-DecrementMemberCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Project-DecrementMemberCount", err, "go_routine")
		}
	} else {
		project.NumberOfMembers--
		if err := initializers.DB.Save(&project).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Project-DecrementMemberCount", err, "go_routine")
		}
		setUserCollaborativeProject(project.UserID)
	}
}

func IncrementUserProject(userID uuid.UUID) {
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No User of this ID found-IncrementEventCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching User-IncrementEventCount", err, "go_routine")
		}
	} else {
		user.NoOfProjects++
		if err := initializers.DB.Save(&user).Error; err != nil {
			helpers.LogDatabaseError("Error while updating User-IncrementEventCount", err, "go_routine")
		}
	}
}

func DecrementUserProject(userID uuid.UUID) {
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No User of this ID found-DecrementEventCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching User-DecrementEventCount", err, "go_routine")
		}
	} else {
		user.NoOfProjects--
		if err := initializers.DB.Save(&user).Error; err != nil {
			helpers.LogDatabaseError("Error while updating User-DecrementEventCount", err, "go_routine")
		}
	}
}

func setUserCollaborativeProject(userID uuid.UUID) {
	var user models.User
	if err := initializers.DB.Preload("Projects").First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No User of this ID found-DecrementEventCount.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching User-DecrementEventCount", err, "go_routine")
		}
	} else {
		user.NoOfCollaborativeProjects = 0
		for _, project := range user.Projects {
			if project.NumberOfMembers > 1 {
				user.NoOfCollaborativeProjects++
			}
		}
		if err := initializers.DB.Save(&user).Error; err != nil {
			helpers.LogDatabaseError("Error while updating User-DecrementEventCount", err, "go_routine")
		}
	}
}

func IncrementReposts(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post of this ID found-IncrementReposts.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post-IncrementReposts", err, "go_routine")
		}
	} else {
		post.NoOfReposts++
		if err := initializers.DB.Save(&post).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Post-IncrementReposts", err, "go_routine")
		}
	}
}

func IncrementResourceBucketFiles(resourceBucketID uuid.UUID) {
	var resourceBucket models.ResourceBucket
	if err := initializers.DB.First(&resourceBucket, "id = ?", resourceBucketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Resource Bucket of this ID found-IncrementResourceBucketFiles.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	} else {
		resourceBucket.NumberOfFiles++
		if err := initializers.DB.Save(&resourceBucket).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	}
}

func DecrementResourceBucketFiles(resourceBucketID uuid.UUID) {
	var resourceBucket models.ResourceBucket
	if err := initializers.DB.First(&resourceBucket, "id = ?", resourceBucketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Resource Bucket of this ID found-IncrementResourceBucketFiles.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	} else {
		resourceBucket.NumberOfFiles--
		if err := initializers.DB.Save(&resourceBucket).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	}
}
