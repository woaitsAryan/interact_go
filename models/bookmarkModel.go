package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectBookmark struct {
	ID        uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID             `gorm:"type:uuid;not null" json:"userId"`
	Title     string                `json:"title"`
	CreatedAt time.Time             `json:"createdAt"`
	Items     []ProjectBookmarkItem `gorm:"foreignKey:ProjectBookmarkID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

type ProjectBookmarkItem struct {
	ProjectBookmarkID uuid.UUID `gorm:"type:uuid;not null"`
	ProjectID         uuid.UUID `gorm:"type:uuid;not null"`
}

type PostBookmark struct {
	ID        uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null" json:"userId"`
	Title     string             `json:"title"`
	CreatedAt time.Time          `json:"createdAt"`
	Items     []PostBookmarkItem `gorm:"foreignKey:PostBookmarkID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

type PostBookmarkItem struct {
	PostBookmarkID uuid.UUID `gorm:"type:uuid;not null"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null"`
}
