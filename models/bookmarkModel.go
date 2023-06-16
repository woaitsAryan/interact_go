package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectBookmark struct {
	ID        uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	Title     string                `json:"title"`
	CreatedAt time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	Items     []ProjectBookmarkItem `gorm:"foreignKey:ProjectBookmarkID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

type ProjectBookmarkItem struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectBookmarkID uuid.UUID `gorm:"type:uuid;not null"`
	ProjectID         uuid.UUID `gorm:"type:uuid;not null"`
	Project           Project
}

type PostBookmark struct {
	ID        uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null" json:"userID"`
	Title     string             `json:"title"`
	CreatedAt time.Time          `gorm:"default:current_timestamp" json:"createdAt"`
	Items     []PostBookmarkItem `gorm:"foreignKey:PostBookmarkID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

type PostBookmarkItem struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostBookmarkID uuid.UUID `gorm:"type:uuid;not null"`
	PostID         uuid.UUID `gorm:"type:uuid;not null"`
	Post           Post
}
