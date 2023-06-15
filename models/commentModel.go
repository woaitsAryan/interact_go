package models

import (
	"time"

	"github.com/google/uuid"
)

type PostComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"postId"`
	Post      Post      `gorm:"" json:"post"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"commentedById"`
	User      User      `gorm:"" json:"commentedBy"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProjectComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectId"`
	Project   Project   `gorm:"" json:"project"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"commentedById"`
	User      User      `gorm:"" json:"commentedBy"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
