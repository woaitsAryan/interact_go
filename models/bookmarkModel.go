package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectBookmark struct {
	gorm.Model
	UserID uuid.UUID
	Title  string `json:"title"`
	// Projects  []Project `gorm:"foreignKey:BookmarkID"`
}

type PostBookmark struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"postId"`
	Post      Post      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
