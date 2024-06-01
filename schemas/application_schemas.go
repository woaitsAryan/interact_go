package schemas

import "github.com/lib/pq"

type ApplicationCreateSchema struct {
	Content       string         `json:"content" validate:"max=500"`
	Links         pq.StringArray `json:"links" validate:"dive,url"`
	IncludeEmail  bool           `json:"includeEmail"`
	IncludeResume bool           `json:"includeResume"`
	YOE           int8           `json:"yoe"`
}
