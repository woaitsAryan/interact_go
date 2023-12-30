package cache

import (
	"github.com/Pratham-Mishra04/interact/models"
)

func GetEvent(slug string) (*models.Event, error) {
	var event models.Event
	err := GetFromCacheGeneric("event-"+slug, &event)
	return &event, err
}

func SetEvent(slug string, event *models.Event) error {
	return SetToCacheGeneric("event-"+slug, event)
}

func RemoveEvent(slug string) error {
	return RemoveFromCacheGeneric("event-" + slug)
}
