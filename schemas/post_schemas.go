package schemas

import "github.com/lib/pq"

type PostCreateSchema struct { // from request
	Content         string         `json:"content" validate:"required,max=2000"`
	Tags            pq.StringArray `json:"tags" validate:"dive,alphanum"`
	TaggedUsernames pq.StringArray `json:"taggedUsernames"`
	RePostID        string         `json:"rePostID"`
}

type PostUpdateSchema struct {
	Content         string          `json:"content" validate:"max=2000"`
	Tags            *pq.StringArray `json:"tags" validate:"dive,alphanum"`
	TaggedUsernames pq.StringArray  `json:"taggedUsernames"`
}
