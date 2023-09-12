package schemas

type ApplicationCreateSchema struct {
	Content string   `json:"content" validate:"max=500"`
	Links   []string `json:"links" validate:"url,dive"`
}
