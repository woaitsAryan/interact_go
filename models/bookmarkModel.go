package models

import (
	"time"

	"github.com/google/uuid"
)

type PostBookmark struct {
	ID        uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null" json:"userID"`
	Title     string             `json:"title"`
	CreatedAt time.Time          `gorm:"default:current_timestamp" json:"createdAt"`
	PostItems []PostBookmarkItem `gorm:"foreignKey:PostBookmarkID;constraint:OnDelete:CASCADE" json:"postItems,omitempty"`
}

type PostBookmarkItem struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostBookmarkID uuid.UUID `gorm:"type:uuid;not null" json:"postBookmarkID"`
	PostID         uuid.UUID `gorm:"type:uuid;not null" json:"postID"`
	Post           Post      `json:"post"`
}
type ProjectBookmark struct {
	ID           uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID       uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	Title        string                `json:"title"`
	CreatedAt    time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	ProjectItems []ProjectBookmarkItem `gorm:"foreignKey:ProjectBookmarkID;constraint:OnDelete:CASCADE" json:"projectItems,omitempty"`
}

type ProjectBookmarkItem struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectBookmarkID uuid.UUID `gorm:"type:uuid;not null" json:"projectBookmarkID"`
	ProjectID         uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Project           Project   `json:"project"`
}
