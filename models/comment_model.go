package models

import (
	"time"

	"github.com/google/uuid"
)

type PostComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid" json:"postID"`
	Post      Post      `gorm:"" json:"post"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"" json:"user"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	NoLikes   int       `json:"noLikes"`
	Edited    bool      `gorm:"default:false" json:"edited"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
type UserPostCommentLike struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null" json:"userID"`
	User          User        `gorm:"" json:"user"`
	PostCommentID uuid.UUID   `gorm:"type:uuid;not null" json:"commentID"`
	PostComment   PostComment `gorm:"" json:"comment"`
	CreatedAt     time.Time   `gorm:"default:current_timestamp" json:"createdAt"`
}

type ProjectComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid" json:"projectID"`
	Project   Project   `gorm:"" json:"project"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"" json:"user"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	NoLikes   int       `json:"noLikes"`
	Edited    bool      `gorm:"default:false" json:"edited"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
type UserProjectCommentLike struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"userID"`
	User             User           `gorm:"" json:"user"`
	ProjectCommentID uuid.UUID      `gorm:"type:uuid;not null" json:"commentID"`
	ProjectComment   ProjectComment `gorm:"" json:"comment"`
	CreatedAt        time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
}
