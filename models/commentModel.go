package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectId"`
	Project   Project   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"postId"`
	Post      Post      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
