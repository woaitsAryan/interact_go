package schemas

/*
	Request body for creating a poll

used in controllers/organization_controllers/reviews.go
*/
type ReviewCreateSchema struct {
	Content   string `json:"content" validate:"required,max=500"`
	Rating    int    `json:"rating"  validate:"required,min=1,max=5"`
	Anonymous bool   `json:"isAnonymous"`
}
