package routines

import (
	"sync"
	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"gorm.io/gorm"
)

//* Posts

func IncrementPostImpression(posts []models.Post) {
	workerCount := 5
	var postIDs []string

	for _, post := range posts {
		key := post.ID.String()
		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return
		} else if impressionCount >= 9 {
			postIDs = append(postIDs, key)
			cache.ResetImpression(key)
		} else {
			cache.IncrementImpression(key)
		}
	}
	incrementPostImpressionsConcurrently(postIDs, workerCount)
}

func incrementPostImpressions(postID string, ch chan<- uint) {
	result := initializers.DB.Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", 10))
	if result.Error != nil {
		helpers.LogDatabaseError("Error updating PostImpressionCount: %v", result.Error, "impression_routines.go")
		return
	}
	ch <- 1
}

func incrementPostImpressionsConcurrently(postIDs []string, workerCount int) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup
	for _, postID := range postIDs {
		wg.Add(1)
		go func(postID string) {
			defer wg.Done()

			incrementPostImpressions(postID, done)
		}(postID)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
}

//* Projects

func IncrementProjectImpression(projects []models.Project) {
	workerCount := 5
	var projectIDs []string

	for _, project := range projects {
		key := project.ID.String()
		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return
		} else if impressionCount >= 9 {
			projectIDs = append(projectIDs, key)
			cache.ResetImpression(key)
		} else {
			cache.IncrementImpression(key)
		}
	}
	incrementProjectImpressionsConcurrently(projectIDs, workerCount)
}

func incrementProjectImpressions(projectID string, ch chan<- uint) {
	result := initializers.DB.Model(&models.Project{}).Where("id = ?", projectID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", 10))
	if result.Error != nil {
		helpers.LogDatabaseError("Error updating ProjectImpressionCount: %v", result.Error, "impression_routines.go")
		return
	}
	ch <- 1
}

func incrementProjectImpressionsConcurrently(projectIDs []string, workerCount int) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup
	for _, projectID := range projectIDs {
		wg.Add(1)
		go func(projectID string) {
			defer wg.Done()

			incrementProjectImpressions(projectID, done)
		}(projectID)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
}

//* Events

func IncrementEventImpression(events []models.Event) {
	workerCount := 5
	var eventIDs []string

	for _, event := range events {
		key := event.ID.String()
		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return
		} else if impressionCount >= 9 {
			eventIDs = append(eventIDs, key)
			cache.ResetImpression(key)
		} else {
			cache.IncrementImpression(key)
		}
	}
	incrementEventImpressionsConcurrently(eventIDs, workerCount)
}

func incrementEventImpressions(eventID string, ch chan<- uint) {
	result := initializers.DB.Model(&models.Event{}).Where("id = ?", eventID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", 10))
	if result.Error != nil {
		helpers.LogDatabaseError("Error updating EventImpressionCount: %v", result.Error, "impression_routines.go")
		return
	}
	ch <- 1
}

func incrementEventImpressionsConcurrently(eventIDs []string, workerCount int) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup
	for _, eventID := range eventIDs {
		wg.Add(1)
		go func(eventID string) {
			defer wg.Done()
			incrementEventImpressions(eventID, done)
		}(eventID)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
}

//* Openings

func IncrementOpeningImpression(openings []models.Opening) {
	workerCount := 5
	var openingIDs []string

	for _, opening := range openings {
		key := opening.ID.String()
		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return
		} else if impressionCount >= 9 {
			openingIDs = append(openingIDs, key)
			cache.ResetImpression(key)
		} else {
			cache.IncrementImpression(key)
		}
	}
	incrementOpeningImpressionsConcurrently(openingIDs, workerCount)
}

func incrementOpeningImpressions(openingID string, ch chan<- uint) {
	result := initializers.DB.Model(&models.Opening{}).Where("id = ?", openingID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", 10))
	if result.Error != nil {
		helpers.LogDatabaseError("Error updating OpeningImpressionCount: %v", result.Error, "impression_routines.go")
		return
	}
	ch <- 1
}

func incrementOpeningImpressionsConcurrently(openingIDs []string, workerCount int) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup
	for _, openingID := range openingIDs {
		wg.Add(1)
		go func(openingID string) {
			defer wg.Done()
			incrementOpeningImpressions(openingID, done)
		}(openingID)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
}

//* User

func IncrementUserImpression(users []models.User) {
	workerCount := 5
	var userIDs []string

	for _, user := range users {
		key := user.ID.String()
		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return
		} else if impressionCount >= 9 {
			userIDs = append(userIDs, key)
			cache.ResetImpression(key)
		} else {
			cache.IncrementImpression(key)
		}
	}
	incrementUserImpressionsConcurrently(userIDs, workerCount)
}

func incrementUserImpressions(userID string, ch chan<- uint) {
	result := initializers.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", 10))
	if result.Error != nil {
		helpers.LogDatabaseError("Error updating UserImpressionCount: %v", result.Error, "impression_routines.go")
		return
	}
	ch <- 1
}

func incrementUserImpressionsConcurrently(userIDs []string, workerCount int) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup
	for _, userID := range userIDs {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			incrementUserImpressions(userID, done)
		}(userID)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
}