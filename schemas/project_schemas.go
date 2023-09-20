package schemas

import "github.com/lib/pq"

type ProjectCreateSchema struct {
	Title       string         `json:"title" validate:"required,max=20"`
	Tagline     string         `json:"tagline" validate:"required"`
	Description string         `json:"description" validate:"required,max=500"`
	Tags        pq.StringArray `json:"tags" validate:"dive,alpha"`
	Category    string         `json:"category" validate:"required"`
	IsPrivate   bool           `json:"isPrivate" validate:"boolean"`
	Links       pq.StringArray `json:"links" validate:"dive,url"`
}

type ProjectUpdateSchema struct {
	Tagline      string         `json:"tagline" validate:"alphanum,max=40"`
	CoverPic     string         `json:"coverPic" validate:"image"`
	Description  string         `json:"description" validate:"alphanum,max=500"`
	Page         string         `json:"page"`
	Tags         pq.StringArray `json:"tags" validate:"alpha,dive"`
	IsPrivate    bool           `json:"isPrivate" validate:"boolean"`
	Links        pq.StringArray `json:"links" validate:"url,dive"`
	PrivateLinks pq.StringArray `json:"privateLinks" validate:"url,dive"`
}
