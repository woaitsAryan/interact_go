package schemas

/*
	Request body for creating a poll

Options field is just an array of options as strings
*/
type CreatePollRequest struct {
	Title         string   `json:"title" validate:"max=50"`
	Content       string   `json:"content" validate:"required,max=500"`
	Options       []string `json:"options"`
	IsMultiAnswer bool     `json:"isMultiAnswer"`
	IsOpen        bool     `json:"isOpen"`
}

type EditPollRequest struct {
	Content string `json:"question" validate:"max=500"`
	IsOpen  bool   `json:"isOpen"`
}
