package schemas

import "github.com/lib/pq"

type PostCreateScheam struct {
	Content string         `json:"content" validate:"required,max=500"`
	Images  pq.StringArray `json:"images"`
	Tags    pq.StringArray `json:"tags"`
}

type PostUpdateScheam struct {
	Content string         `json:"content" validate:"max=500"`
	Tags    pq.StringArray `json:"tags"`
}
