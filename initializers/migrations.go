package initializers

import "github.com/Pratham-Mishra04/interact/models"

func AutoMigrate() {

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
		&models.Message{},
		&models.ProjectChatMessage{},
		&models.Chat{},
		&models.ProjectChat{},
		&models.Post{},
		&models.Project{},
		&models.ProjectView{},
		&models.Opening{},
		&models.Application{},
		&models.Membership{},
		&models.PostComment{},
		&models.ProjectComment{},
		&models.ChatInvitation{},
		&models.ProjectInvitation{},
	)

}
