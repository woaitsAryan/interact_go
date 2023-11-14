package initializers

import (
	"fmt"

	"github.com/Pratham-Mishra04/interact/models"
)

func AutoMigrate() {
	fmt.Println("\nStarting Migrations...")
	DB.AutoMigrate(
		&models.User{},
		&models.Achievement{},
		&models.ProfileView{},
		&models.FollowFollower{},
		&models.Notification{},
		&models.PostBookmark{},
		&models.PostBookmarkItem{},
		&models.ProjectBookmark{},
		&models.ProjectBookmarkItem{},
		&models.OpeningBookmark{},
		&models.OpeningBookmarkItem{},

		&models.Message{},
		&models.GroupChatMessage{},
		&models.Chat{},
		&models.GroupChat{},
		&models.GroupChatMembership{},

		&models.Post{},
		&models.Project{},
		&models.ProjectView{},
		&models.ProjectHistory{},
		&models.Task{},
		&models.SubTask{},
		&models.Opening{},
		&models.Application{},
		&models.Membership{},
		&models.Comment{},
		&models.Invitation{},
		&models.Like{},
		&models.UserPostTag{},
		&models.LastViewedProjects{},
		&models.LastViewedOpenings{},
		&models.SearchQuery{},
		&models.UserVerification{},
		&models.OAuth{},
		&models.Organization{},
		&models.OrganizationMembership{},
		&models.OrganizationHistory{},
		&models.Report{},
	)
	fmt.Println("Migrations Finished!")
}
