package cache

import (
	"github.com/Pratham-Mishra04/interact/models"
)

func GetOrganization(slug string) (*models.Organization, error) {
	var organization models.Organization
	err := GetFromCacheGeneric("organization-"+slug, &organization)
	return &organization, err
}

func SetOrganization(slug string, organization *models.Organization) error {
	return SetToCacheGeneric("organization-"+slug, organization)
}

func RemoveOrganization(slug string) error {
	return RemoveFromCacheGeneric("organization-" + slug)
}

func GetResourceBucket(slug string) (*models.ResourceBucket, error) {
	var resourceBucket models.ResourceBucket
	err := GetFromCacheGeneric("resource_bucket-"+slug, &resourceBucket)
	return &resourceBucket, err
}

func SetResourceBucket(slug string, resourceBucket *models.ResourceBucket) error {
	return SetToCacheGeneric("resource_bucket-"+slug, resourceBucket)
}

func RemoveResourceBucket(slug string) error {
	return RemoveFromCacheGeneric("resource_bucket-" + slug)
}
