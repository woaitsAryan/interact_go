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
	PostItems []PostBookmarkItem `gorm:"foreignKey:PostBookmarkID;constraint:OnDelete:CASCADE" json:"postItems"`
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
	ProjectItems []ProjectBookmarkItem `gorm:"foreignKey:ProjectBookmarkID;constraint:OnDelete:CASCADE" json:"projectItems"`
}

type ProjectBookmarkItem struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectBookmarkID uuid.UUID `gorm:"type:uuid;not null" json:"projectBookmarkID"`
	ProjectID         uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Project           Project   `json:"project"`
}

type OpeningBookmark struct {
	ID           uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID       uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	Title        string                `json:"title"`
	CreatedAt    time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	OpeningItems []OpeningBookmarkItem `gorm:"foreignKey:OpeningBookmarkID;constraint:OnDelete:CASCADE" json:"openingItems"`
}

type OpeningBookmarkItem struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OpeningBookmarkID uuid.UUID `gorm:"type:uuid;not null" json:"openingBookmarkID"`
	OpeningID         uuid.UUID `gorm:"type:uuid;not null" json:"openingID"`
	Opening           Opening   `json:"opening"`
}
