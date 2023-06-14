package schemas

import (
	"github.com/lib/pq"
)

type UserCreateSchema struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	PhoneNo         string `json:"phoneNo"`
	ProfilePic      string `json:"profilePic"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type UserUpdateSchema struct {
	Name       string         `json:"name"`
	PhoneNo    string         `json:"phoneNo"`
	ProfilePic string         `json:"profilePic"`
	CoverPic   string         `json:"coverPic"`
	Bio        string         `json:"bio"`
	Title      string         `json:"title"`
	Tagline    string         `json:"tagline"`
	Tags       pq.StringArray `gorm:"type:text[]" json:"tags"`
}
