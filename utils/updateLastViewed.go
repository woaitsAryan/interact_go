package utils

import (
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateLastViewed(userID uuid.UUID, projectID uuid.UUID) error {

	var projectView models.LastViewed
	if err := initializers.DB.Preload("User").Where("user_id = ? AND project_id=?", userID, projectID).First(&projectView).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			projectView.ProjectID = projectID
			projectView.UserID = userID
			projectView.Timestamp = time.Now()

			result := initializers.DB.Create(&projectView)
			if result.Error != nil {
				return &fiber.Error{Code: 500, Message: "Database Error whiling creating last viewed."}
			}
		} else {
			return &fiber.Error{Code: 500, Message: "Database Error."}
		}
	} else {
		projectView.Timestamp = time.Now()

		result := initializers.DB.Save(&projectView)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Database Error whiling updating last viewed."}
		}
	}

	return nil

}
