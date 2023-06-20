package schemas

import "github.com/lib/pq"

type ApplicationCreateScheam struct {
	Content string         `json:"content" validate:"max=500"`
	Links   pq.StringArray `json:"links"`
}
