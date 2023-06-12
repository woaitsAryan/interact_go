package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name              string
	Username          string
	Email             string `gorm:"unique"`
	Password          string
	RegNo             string `gorm:"unique"`
	ProfilePic        string
	PhoneNo           string
	PasswordChangedAt time.Time `gorm:"default:current_timestamp"`
	Admin             bool      `gorm:"default:false"`
	Active            bool      `gorm:"default:true"`
}

type UserCreateSchema struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// This func is called before gorm conversion
// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.PasswordChangedAt = time.Now()
// 	return nil
// }
