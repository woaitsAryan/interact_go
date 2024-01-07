package utils

import (
	"sort"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

type ViewResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func GetProfileViews(userID uuid.UUID) ([]ViewResponse, int, error) {
	// Create a map to store the count for each date
	viewsMap := make(map[time.Time]int)

	// Initialize the map with all past 7 dates and count as 0
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).UTC().Truncate(24 * time.Hour)
		viewsMap[date] = 0
	}

	// Retrieve the profile views from the database for the past 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -6).UTC().Truncate(24 * time.Hour) // Get the date 7 days ago
	var profileViews []models.ProfileView
	if err := initializers.DB.Where("user_id = ? AND date >= ?", userID, sevenDaysAgo).Find(&profileViews).Error; err != nil {
		return nil, 0, helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// Update the count in the map based on the retrieved profile views
	for _, view := range profileViews {
		date := view.Date.Truncate(24 * time.Hour)
		viewsMap[date] += view.Count
	}

	// Convert the map entries to ViewResponse objects
	var response []ViewResponse
	var totalViews int
	for date := range viewsMap {
		response = append(response, ViewResponse{
			Date:  date,
			Count: viewsMap[date],
		})
		totalViews += viewsMap[date]
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Date.Before(response[j].Date)
	})

	return response, totalViews, nil
}

func GetProjectViews(projectID uuid.UUID) ([]ViewResponse, int, error) {

	var projectViews []models.ProjectView
	if err := initializers.DB.Where("project_id = ? AND date >= ?", projectID, time.Now().AddDate(0, 0, -30).Format("2006-01-02")).Find(&projectViews).Error; err != nil {
		return nil, 0, helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var totalViews int

	var response []ViewResponse
	for _, view := range projectViews {
		response = append(response, ViewResponse{
			Date:  view.Date,
			Count: view.Count,
		})
		totalViews += view.Count
	}

	return response, totalViews, nil
}
