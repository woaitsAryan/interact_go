package schemas

import (
	"github.com/lib/pq"
)

type EventCreateSchema struct {
	Title       string         `json:"title" validate:"required,max=25"`
	Tagline     string         `json:"tagline" validate:"required,max=50"`
	Description string         `json:"description" validate:"required,max=1000"`
	Tags        pq.StringArray `json:"tags"`
	Category    string         `json:"category" validate:"required"`
	Location    string         `json:"location" validate:"required"`
	Links       pq.StringArray `json:"links" validate:"dive,url"`
	StartTime   string         `json:"startTime"`
	EndTime     string         `json:"endTime"`
}

type EventUpdateSchema struct {
	Tagline     string         `json:"tagline" validate:"max=50"`
	CoverPic    string         `json:"coverPic" validate:"image"`
	Category    string         `json:"category"`
	Description string         `json:"description" validate:"max=1000"`
	Location    string         `json:"location"`
	Tags        pq.StringArray `json:"tags"`
	Links       pq.StringArray `json:"links" validate:"dive,url"`
	StartTime   string         `json:"startTime"`
	EndTime     string         `json:"endTime"`
}
