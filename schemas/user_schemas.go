package schemas

import (
	"github.com/lib/pq"
)

type UserCreateSchema struct {
	Name            string `json:"name" validate:"required,max=25"`
	Username        string `json:"username" validate:"required,max=16"` //alphanum+_
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type UserUpdateSchema struct {
	Name       *string         `json:"name" validate:"max=25"`
	ProfilePic *string         `json:"profilePic" validate:"image"`
	CoverPic   *string         `json:"coverPic" validate:"image"`
	Bio        *string         `json:"bio" validate:"max=500"`
	Title      *string         `json:"title"`
	Tagline    *string         `json:"tagline" validate:"max=25"`
	Tags       *pq.StringArray `json:"tags"`
	Links      *pq.StringArray `json:"links"`
}

type ProfileUpdateSchema struct {
	School      *string         `json:"school" validate:"max=25"`
	Degree      *string         `json:"degree" validate:"max=25"`
	YOG         *string         `json:"yog"`
	Description *string         `json:"description" validate:"max=1500"`
	Hobbies     *pq.StringArray `json:"hobbies"`
	Areas       *pq.StringArray `json:"areas"`
	Email       *string         `json:"email"`
	PhoneNo     *string         `json:"phoneNo"`
	Location    *string         `json:"location"`
}

type AchievementCreateSchema struct {
	Achievements []AchievementSchema `json:"achievements"`
}

type AchievementSchema struct {
	ID     string   `json:"id"`
	Title  string   `json:"title" validate:"alpha"`
	Skills []string `json:"skills" validate:"dive,alpha"`
}
