package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectBookmark struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID    uuid.UUID
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProjectBookmarkItem struct {
	ProjectBookmarkID uuid.UUID
	ProjectID         uuid.UUID
}

type PostBookmark struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostBookmarkItem struct {
	PostBookmarkID uuid.UUID
	ProjectID      uuid.UUID
}
