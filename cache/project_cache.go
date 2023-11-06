package cache

import (
	"encoding/json"
	"fmt"

	"github.com/Pratham-Mishra04/interact/models"
)

func GetProject(id string) (*models.Project, error) {
	data, err := GetFromCache("project-" + id)
	if err != nil {
		return nil, err
	}

	project := models.Project{}
	if err = json.Unmarshal([]byte(data), &project); err != nil {
		return nil, fmt.Errorf("error while unmarshaling project: %w", err)
	}
	return &project, nil
}

func SetProject(id string, project *models.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("error while marshaling project: %w", err)
	}
	if err := SetToCache("project-"+id, data); err != nil {
		return err
	}
	return nil
}
