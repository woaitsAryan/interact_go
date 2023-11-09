package schemas

import "github.com/lib/pq"

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
	Bio        *string         `json:"bio" validate:"max=300"`
	Title      *string         `json:"title"`
	Tagline    *string         `json:"tagline" validate:"max=25"`
	Tags       *pq.StringArray `json:"tags" validate:"dive,alpha"`
	Links      *pq.StringArray `json:"links" validate:"dive,url"`
}

type AchievementCreateSchema struct {
	Achievements []AchievementSchema `json:"achievements"`
}

type AchievementSchema struct {
	ID     string   `json:"id"`
	Title  string   `json:"title" validate:"alpha"`
	Skills []string `json:"skills" validate:"dive,alpha"`
}
