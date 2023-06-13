package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name                      string         `gorm:"varchar(25);not null" json:"name"`
	Username                  string         `gorm:"varchar(10);unique;not null" json:"username"`
	Email                     string         `gorm:"unique;not null" json:"email"`
	Password                  string         `json:"password"`
	ProfilePic                string         `json:"profilePic"`
	CoverPic                  string         `json:"coverPic"`
	PhoneNo                   string         `json:"phoneNo"`
	Bio                       string         `json:"bio"`
	Title                     string         `json:"title"`
	Tagline                   string         `json:"tagline"`
	Tags                      pq.StringArray `gorm:"type:text[]" json:"tags"`
	Followers                 []*User        `gorm:"many2many:user_followers;joinForeignKey:follower_id;joinReferences:id" json:"followers,omitempty"`
	Following                 []*User        `gorm:"many2many:user_followers;joinForeignKey:user_id;joinReferences:id" json:"following,omitempty"`
	LastViewed                []*Project     `gorm:"many2many:user_last_viewed_projects;" json:"lastViewed,omitempty"`
	PasswordResetToken        string         `json:"-"`
	PasswordResetTokenExpires time.Time      `json:"-"`
	PasswordChangedAt         time.Time      `gorm:"default:current_timestamp" json:"-"`
	Admin                     bool           `gorm:"default:false" json:"admin"`
	Active                    bool           `gorm:"default:true" json:"active"`
}

type ProfileView struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Date   time.Time `json:"date"`
	Count  int       `json:"count"`
}

type UserCreateSchema struct {
	Name       string `json:"name" validate:"required"`
	Username   string `json:"username" validate:"required"`
	PhoneNo    string `json:"phoneNo"`
	ProfilePic string `json:"profilePic"`

	Email             string    `json:"email" validate:"required,email"`
	Password          string    `json:"password" validate:"required,min=8"`
	ConfirmPassword   string    `json:"confirmPassword" validate:"required,min=8"`
	PasswordChangedAt time.Time `json:"-"`
	Admin             bool      `json:"-"`
	Active            bool      `json:"-"`
}

// This func is called before gorm conversion
// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.PasswordChangedAt = time.Now()
// 	return nil
// }
