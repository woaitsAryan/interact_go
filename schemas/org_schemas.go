package schemas

import (
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/lib/pq"
)

/*
	Request body for creating a poll

used in controllers/organization_controllers/reviews.go
*/
type ReviewCreateSchema struct {
	Content   string `json:"content" validate:"required,max=500"`
	Rating    int8   `json:"rating"  validate:"required,min=1,max=5"`
	Anonymous bool   `json:"isAnonymous"`
}

type ResourceBucketCreateSchema struct {
	Title       string                  `json:"title" validate:"required,max=50"`
	Description string                  `json:"description" validate:"max=500"`
	ViewAccess  models.OrganizationRole `json:"viewAccess" validate:"required"`
	EditAccess  models.OrganizationRole `json:"editAccess" validate:"required"`
}

type ResourceBucketEditSchema struct {
	Title       string                  `json:"title" validate:"max=50"`
	Description string                  `json:"description" validate:"max=500"`
	ViewAccess  models.OrganizationRole `json:"viewAccess"`
	EditAccess  models.OrganizationRole `json:"editAccess"`
}

type ResourceFileCreateSchema struct {
	Title       string `json:"title" validate:"required,max=50"`
	Description string `json:"description" validate:"max=500"`
}

type ResourceFileEditSchema struct {
	Title       string `json:"title" validate:"max=50"`
	Description string `json:"description" validate:"max=500"`
}

type AnnouncementCreateSchema struct {
	Title           string         `json:"title" validate:"max=50"`
	Content         string         `json:"content" validate:"required,max=1000"`
	TaggedUsernames pq.StringArray `json:"taggedUsernames"`
	IsOpen          bool           `json:"isOpen"`
}

type AnnouncementUpdateSchema struct {
	Title   string `json:"title" validate:"max=50"`
	Content string `json:"content" validate:"max=1000"`
	IsOpen  bool   `json:"isOpen"`
}
