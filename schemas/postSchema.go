package schemas

import "github.com/lib/pq"

type PostCreateScheam struct { // from request
	Content string `json:"content" validate:"required,max=500"`
	// Images  pq.StringArray `json:"images" validate:"dive,image"`
	Tags pq.StringArray `json:"tags" validate:"dive,alphanum"`
}

type PostUpdateScheam struct {
	Content string         `json:"content" validate:"max=500"`
	Tags    pq.StringArray `json:"tags" validate:"dive,alphanum"`
}
