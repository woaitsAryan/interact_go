package initializers

import "github.com/Pratham-Mishra04/interact/models"

func AutoMigrate() {

	DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Project{},
		&models.Application{},
		&models.Chat{},
		&models.Membership{},
		&models.Comment{},
		&models.Message{},
		&models.Notification{},
		&models.Opening{},
		&models.PostBookmark{},
		&models.ProfileView{},
		&models.ProjectBookmark{},
		&models.ProjectView{},
		&models.ProjectBookmarkItem{},
		&models.PostBookmarkItem{},
	)
}
