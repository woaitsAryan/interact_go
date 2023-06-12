package initializers

import "github.com/Pratham-Mishra04/interact/models"

func AutoMigrate() {
	DB.AutoMigrate(&models.User{})
}
