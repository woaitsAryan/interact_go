package initializers

import "github.com/Pratham-Mishra04/interact/models"

func AutoMigrate() {

	// DB.AutoMigrate(&models.User{}, &models.Post{})

	// // DB.AutoMigrate(&models.User{})
	// // DB.AutoMigrate(&models.Post{})
	// // DB.AutoMigrate(&models.Project{})
	// // DB.AutoMigrate(&models.Application{})
	// // DB.AutoMigrate(&models.Chat{})
	// // DB.AutoMigrate(&models.Collaborator{})
	// // DB.AutoMigrate(&models.Invitation{})
	// // DB.AutoMigrate(&models.Comment{})
	// // DB.AutoMigrate(&models.Message{})
	// // DB.AutoMigrate(&models.Notification{})
	// // DB.AutoMigrate(&models.Opening{})
	// // DB.AutoMigrate(&models.PostBookmark{})
	// // DB.AutoMigrate(&models.ProfileView{})
	// DB.AutoMigrate(&models.ProjectBookmark{})

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
	)
}
