package schemas

import (
	"github.com/lib/pq"
)

type UserCreateSchema struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"alphanum,required"`
	PhoneNo         string `json:"phoneNo" validate:"e164"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type UserUpdateSchema struct {
	Name       string         `json:"name" validate:"alpha"`
	PhoneNo    string         `json:"phoneNo"  validate:"e164"`
	ProfilePic string         `json:"profilePic" validate:"image"`
	CoverPic   string         `json:"coverPic" validate:"image"`
	Bio        string         `json:"bio"`
	Title      string         `json:"title"`
	Tagline    string         `json:"tagline"`
	Tags       pq.StringArray `json:"tags" validate:"dive,alpha"`
}
