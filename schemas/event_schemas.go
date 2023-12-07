package schemas

import (
	"time"

	"github.com/lib/pq"
)

type EventCreateSchema struct {
	Title       string         `json:"title" validate:"required,max=20"`
	Tagline     string         `json:"tagline" validate:"required,max=40"`
	Description string         `json:"description" validate:"required,max=1000"`
	Tags        pq.StringArray `json:"tags" validate:"dive"`
	Category    string         `json:"category" validate:"required"`
	Links       pq.StringArray `json:"links" validate:"dive,url"`
	EventDate   time.Time      `json:"date"`
}

type EventUpdateSchema struct {
	Tagline     string         `json:"tagline" validate:"max=40"`
	CoverPic    string         `json:"coverPic" validate:"image"`
	Category    string         `json:"category"`
	Description string         `json:"description" validate:"max=1000"`
	Tags        pq.StringArray `json:"tags" validate:"alpha,dive"`
	Links       pq.StringArray `json:"links" validate:"dive,url"`
	EventDate   time.Time      `json:"date"`
}
