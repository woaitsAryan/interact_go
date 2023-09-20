package schemas

import "github.com/lib/pq"

type PostCreateSchema struct { // from request
	Content       string         `json:"content" validate:"required,max=1000"`
	Tags          pq.StringArray `json:"tags" validate:"dive,alphanum"`
	TaggedUserIDS pq.StringArray `json:"taggedUserIDS"`
	RePostID      string         `json:"rePostID"`
}

type PostUpdateSchema struct {
	Content       string         `json:"content" validate:"max=1000"`
	Tags          pq.StringArray `json:"tags" validate:"dive,alphanum"`
	TaggedUserIDS pq.StringArray `json:"taggedUserIDS"`
}
