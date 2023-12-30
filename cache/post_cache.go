package cache

import (
	"github.com/Pratham-Mishra04/interact/models"
)

func GetPost(slug string) (*models.Post, error) {
	var post models.Post
	err := GetFromCacheGeneric("post-"+slug, &post)
	return &post, err
}

func SetPost(slug string, post *models.Post) error {
	return SetToCacheGeneric("post-"+slug, post)
}

func RemovePost(slug string) error {
	return RemoveFromCacheGeneric("post-" + slug)
}
