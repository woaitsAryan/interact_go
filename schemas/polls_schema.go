package schemas

/* Request body for creating a poll

Options field is just an array of options as strings
*/
type CreatePollRequest struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	IsMultiAnswer bool     `json:"isMultiAnswer"`
}

type EditPollRequest struct {
	Question      string   `json:"question"`
	IsMultiAnswer bool     `json:"isMultiAnswer"`
}