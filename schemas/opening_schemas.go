package schemas

type OpeningCreateSchema struct { // from request
	Title       string   `json:"title" validate:"required,max=50"`
	Description string   `json:"description" validate:"required,max=500"`
	Tags        []string `json:"tags"`
}

type OpeningEditSchema struct { // from request
	Description string   `json:"description" validate:"max=500"`
	Tags        []string `json:"tags"`
	Active      *bool    `json:"active"`
}
