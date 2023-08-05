package schemas

type PostCreateSchema struct { // from request
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"dive,alphanum"`
}

type PostUpdateSchema struct {
	Content string   `json:"content" validate:"max=1000"`
	Tags    []string `json:"tags" validate:"dive,alphanum"`
}
