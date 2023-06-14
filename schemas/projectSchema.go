package schemas

import "github.com/lib/pq"

type ProjectCreateSchema struct {
	Title       string         `json:"title" validate:"alphanum,required,max=20"`
	Tagline     string         `json:"tagline" validate:"alphanum,required,max=40"`
	CoverPic    string         `json:"coverPic" validate:"image"`
	Description string         `json:"description" validate:"alphanum,required,max=500"`
	Tags        pq.StringArray `json:"tags" validate:"alpha, dive"`
	Category    string         `json:"category" validate:"alpha, required"`
	IsPrivate   bool           `json:"isPrivate" validate:"boolean"`
	Links       pq.StringArray `json:"links" validate:"url, dive"`
}

type ProjectUpdateSchema struct {
	Tagline      string         `json:"tagline" validate:"alphanum,max=40"`
	CoverPic     string         `json:"coverPic" validate:"image"`
	Description  string         `json:"description" validate:"alphanum,max=500"`
	Page         string         `json:"page"`
	Tags         pq.StringArray `json:"tags" validate:"alpha, dive"`
	IsPrivate    bool           `json:"isPrivate" validate:"boolean"`
	Links        pq.StringArray `json:"links" validate:"url, dive"`
	PrivateLinks pq.StringArray `json:"privateLinks" validate:"url, dive"`
}
