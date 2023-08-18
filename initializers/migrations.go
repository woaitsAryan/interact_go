package initializers

import (
	"fmt"

	"github.com/Pratham-Mishra04/interact/models"
)

func AutoMigrate() {
	fmt.Println("Starting Migrations...")
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
		// &models.GroupMessage{},
		&models.ProjectChatMessage{},
		&models.Chat{},
		// &models.GroupChat{},
		&models.ProjectChat{},
		&models.ProjectChatMembership{},
		&models.Post{},
		&models.Project{},
		&models.ProjectView{},
		&models.Opening{},
		&models.Application{},
		&models.Membership{},
		&models.Comment{},
		&models.Invitation{},
		&models.UserPostLike{},
		&models.UserProjectLike{},
		&models.UserCommentLike{},
		&models.UserPostTag{},
		&models.LastViewed{},
		&models.SearchQuery{},
		&models.UserVerification{},
		&models.OAuth{},
		&models.Organization{},
		&models.OrganizationMembership{},
		&models.Report{},
	)
	fmt.Println("Migrations Finished!")
}
