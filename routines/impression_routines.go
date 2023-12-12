package routines

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"gorm.io/gorm"
)

type IncrementFunc func(id string, ch chan<- uint)

func IncrementImpressions(items interface{}, getModelID func(interface{}) string, incrementDB IncrementFunc, modelType interface{}) {
	workerCount := 5
	var itemIDs []string

	// Use reflection to get the underlying type of items
	itemsValue := reflect.ValueOf(items)
	if itemsValue.Kind() != reflect.Slice {
		fmt.Printf("Unsupported type: %T\n", items)
		return
	}

	for i := 0; i < itemsValue.Len(); i++ {
		item := itemsValue.Index(i).Interface()
		key := getModelID(item)
		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return
		} else if impressionCount >= 9 {
			itemIDs = append(itemIDs, key)
			cache.ResetImpression(key)
		} else {
			cache.IncrementImpression(key)
		}
	}

	incrementImpressionsConcurrently(itemIDs, workerCount, incrementDB, modelType)
}

func incrementImpressionsConcurrently(ids []string, workerCount int, incrementDB IncrementFunc, modelType interface{}) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			incrementDB(id, done)
		}(id)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
}

func incrementDBImpressions(modelType interface{}, itemID string, ch chan<- uint) {
	result := initializers.DB.Model(modelType).Where("id = ?", itemID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", 10))
	if result.Error != nil {
		typeName := reflect.TypeOf(modelType).Elem().Name()
		helpers.LogDatabaseError(fmt.Sprintf("Error updating %sImpressionCount", typeName), result.Error, "impression_routines")
		return
	}
	ch <- 1
}

// Posts
func IncrementPostImpression(posts []models.Post) {
	IncrementImpressions(posts, func(item interface{}) string {
		return item.(models.Post).ID.String()
	}, incrementDBPostImpressions, &models.Post{})
}

func incrementDBPostImpressions(postID string, ch chan<- uint) {
	incrementDBImpressions(&models.Post{}, postID, ch)
}

// Projects
func IncrementProjectImpression(projects []models.Project) {
	IncrementImpressions(projects, func(item interface{}) string {
		return item.(models.Project).ID.String()
	}, incrementDBProjectImpressions, &models.Project{})
}

func incrementDBProjectImpressions(projectID string, ch chan<- uint) {
	incrementDBImpressions(&models.Project{}, projectID, ch)
}

// Events
func IncrementEventImpression(events []models.Event) {
	IncrementImpressions(events, func(item interface{}) string {
		return item.(models.Event).ID.String()
	}, incrementDBEventImpressions, &models.Event{})
}

func incrementDBEventImpressions(eventID string, ch chan<- uint) {
	incrementDBImpressions(&models.Event{}, eventID, ch)
}

// Openings
func IncrementOpeningImpression(openings []models.Opening) {
	IncrementImpressions(openings, func(item interface{}) string {
		return item.(models.Opening).ID.String()
	}, incrementDBOpeningImpressions, &models.Opening{})
}

func incrementDBOpeningImpressions(openingID string, ch chan<- uint) {
	incrementDBImpressions(&models.Opening{}, openingID, ch)
}

// Users
func IncrementUserImpression(users []models.User) {
	IncrementImpressions(users, func(item interface{}) string {
		return item.(models.User).ID.String()
	}, incrementDBUserImpressions, &models.User{})
}

func incrementDBUserImpressions(userID string, ch chan<- uint) {
	incrementDBImpressions(&models.User{}, userID, ch)
}
