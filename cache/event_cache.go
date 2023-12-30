package cache

import (
	"encoding/json"
	"fmt"

	"github.com/Pratham-Mishra04/interact/models"
)

func GetEvent(slug string) (*models.Event, error) {
	data, err := GetFromCache("event-" + slug)
	if err != nil {
		return nil, err
	}

	event := models.Event{}
	if err = json.Unmarshal([]byte(data), &event); err != nil {
		return nil, fmt.Errorf("error while unmarshaling event: %w", err)
	}
	return &event, nil
}

func SetEvent(slug string, event *models.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error while marshaling event: %w", err)
	}
	if err := SetToCache("event-"+slug, data); err != nil {
		return err
	}
	return nil
}

func RemoveEvent(slug string) error {
	if err := RemoveFromCache("event-" + slug); err != nil {
		return err
	}
	return nil
}
