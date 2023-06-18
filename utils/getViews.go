package utils

import (
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProfileViewResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func GetProfileViews(userID uuid.UUID) ([]ProfileViewResponse, int, error) {

	// Query the database to get the profile views for the past 30 days
	var profileViews []models.ProfileView
	if err := initializers.DB.Where("user_id = ? AND date >= ?", userID, time.Now().AddDate(0, 0, -30).Format("2006-01-02")).Find(&profileViews).Error; err != nil {
		return nil, 0, &fiber.Error{Code: 500, Message: "Failed to get profile views."}
	}

	// Create a slice of ProfileViewResponse objects containing only date and count

	var totalViews int

	var response []ProfileViewResponse
	for _, view := range profileViews {
		response = append(response, ProfileViewResponse{
			Date:  view.Date,
			Count: view.Count,
		})
		totalViews += view.Count
	}

	return response, totalViews, nil
}
