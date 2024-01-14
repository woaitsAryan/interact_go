package schemas

/*
	Request body for creating a poll

Options field is just an array of options as strings
*/
type CreatePollRequest struct {
	Question      string   `json:"question" validate:"required,max=100"`
	Options       []string `json:"options"`
	IsMultiAnswer bool     `json:"isMultiAnswer"`
	IsOpen        bool     `json:"isOpen"`
}

type EditPollRequest struct {
	Question string `json:"question" validate:"required,max=100"`
	IsOpen   bool   `json:"isOpen"`
}
