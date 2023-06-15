package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name                      string         `gorm:"varchar(25);not null" json:"name"`
	Username                  string         `gorm:"varchar(10);unique;not null" json:"username"`
	Email                     string         `gorm:"unique;not null" json:"email"`
	Password                  string         `json:"-"`
	ProfilePic                string         `json:"profilePic"`
	CoverPic                  string         `json:"coverPic"`
	PhoneNo                   string         `json:"phoneNo"`
	Bio                       string         `json:"bio"`
	Title                     string         `json:"title"`
	Tagline                   string         `json:"tagline"`
	Tags                      pq.StringArray `gorm:"type:text[]" json:"tags"`
	PasswordResetToken        string         `json:"-"`
	PasswordResetTokenExpires time.Time      `json:"-"`
	PasswordChangedAt         time.Time      `gorm:"default:current_timestamp" json:"-"`
	Admin                     bool           `gorm:"default:false" json:"-"`
	Active                    bool           `gorm:"default:true" json:"-"`
	CreatedAt                 time.Time      `gorm:"default:current_timestamp" json:"-"`
	Achievements              []Achievement  `gorm:"foreignKey:UserID" json:"achievements,omitempty"`
} //! add last viewed projects

type ProfileView struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Date   time.Time `gorm:"type:date" json:"date"`
	Count  int       `json:"count"`
}

// This func is called before gorm conversion
// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.PasswordChangedAt = time.Now()
// 	return nil
// }

type FollowFollower struct {
	FollowerID uuid.UUID
	Follower   User `gorm:"foreignKey:FollowerID"`
	FollowedID uuid.UUID
	Followed   User      `gorm:"foreignKey:FollowedID"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"-"`
}

type Achievement struct {
	UserID uuid.UUID      `gorm:"type:uuid;not null" json:"userID"`
	User   User           `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	Title  string         `gorm:"type:text;not null" json:"title"`
	Skills pq.StringArray `gorm:"type:text[];not null" json:"skills"`
}
