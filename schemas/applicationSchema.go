package schemas

type ApplicationCreateScheam struct {
	Content string   `json:"content" validate:"max=500"`
	Links   []string `json:"links"`
}
