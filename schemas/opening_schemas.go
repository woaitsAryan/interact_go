package schemas

import "github.com/lib/pq"

type OpeningCreateSchema struct { // from request
	Title       string         `json:"title" validate:"required,max=25"`
	Description string         `json:"description" validate:"required,max=500"`
	Tags        pq.StringArray `json:"tags"`
}

type OpeningEditSchema struct { // from request
	Description string         `json:"description" validate:"max=500"`
	Tags        pq.StringArray `json:"tags"`
	Active      *bool          `json:"active"`
}
