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
