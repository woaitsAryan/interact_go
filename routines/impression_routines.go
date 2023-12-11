package routines

import (
	"sync"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/models"
)

func IncrementPostImpression(posts []models.Post) {
	workerCount := 10

	var ID []string

	for _, post := range posts {
		key := post.ID.String()

		impressionCount, err := cache.GetImpression(key)
		if err != nil {
			return;
		} else if impressionCount >= 9 {
			ID = append(ID, key)
		}else{
			cache.IncrementImpression(key)
		}
	}
	IncrementImpressionsConcurrently(ID, workerCount)
}

func IncrementImpressions(projectID string, ch chan<- uint) {
	cache.IncrementImpression(projectID)
	ch <- 1
}

func IncrementImpressionsConcurrently(IDs []string, workerCount int) {
	done := make(chan uint, workerCount)
	var wg sync.WaitGroup

	for _, ID := range IDs{
		wg.Add(1)

		go func(ID string) {
			defer wg.Done()

			IncrementImpressions(ID, done)
		}(ID)
	}

	go func() {
		wg.Wait()
		close(done)
	}()
}