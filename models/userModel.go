package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name              string    `gorm:"varchar(25);not null" json:"name"`
	Username          string    `gorm:"varchar(10);unique;not null" json:"username"`
	Email             string    `gorm:"unique;not null" json:"email"`
	Password          string    `json:"password"`
	ProfilePic        string    `json:"profilePic"`
	PhoneNo           string    `json:"phoneNo"`
	PasswordChangedAt time.Time `gorm:"default:current_timestamp"`
	Admin             bool      `gorm:"default:false"`
	Active            bool      `gorm:"default:true"`
}

type UserCreateSchema struct {
	Name              string    `json:"name" validate:"required"`
	Username          string    `json:"username" validate:"required"`
	PhoneNo           string    `json:"phoneNo"`
	Email             string    `json:"email" validate:"required,email"`
	Password          string    `json:"password" validate:"required,min=8"`
	PasswordChangedAt time.Time `json:"passwordChangedAt" validate:"-"`
	Admin             bool      `json:"admin" validate:"-"`
	Active            bool      `json:"active" validate:"-"`
}

// This func is called before gorm conversion
// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.PasswordChangedAt = time.Now()
// 	return nil
// }
