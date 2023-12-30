package cache

import (
	"github.com/Pratham-Mishra04/interact/models"
)

func GetProject(slug string) (*models.Project, error) {
	var project models.Project
	err := GetFromCacheGeneric("project-"+slug, &project)
	return &project, err
}

func SetProject(slug string, project *models.Project) error {
	return SetToCacheGeneric("project-"+slug, project)
}

func RemoveProject(slug string) error {
	return RemoveFromCacheGeneric("project-" + slug)
}
