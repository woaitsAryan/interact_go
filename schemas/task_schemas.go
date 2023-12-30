package schemas

import (
	"time"

	"github.com/Pratham-Mishra04/interact/models"
	"github.com/lib/pq"
)

type TaskCreateSchema struct { // from request
	Title       string          `json:"title" validate:"required,max=25"`
	Description string          `json:"description" validate:"max=500"`
	Dateline    time.Time       `json:"deadline" validate:"required"`
	Tags        pq.StringArray  `json:"tags"`
	Users       pq.StringArray  `json:"users" validate:"required"`
	Priority    models.Priority `json:"priority"`
}

type TaskEditSchema struct { // from request
	Title       string          `json:"title" validate:"max=25"`
	Description string          `json:"description" validate:"max=500"`
	Tags        pq.StringArray  `json:"tags"`
	Dateline    time.Time       `json:"deadline"`
	Priority    models.Priority `json:"priority"`
}
