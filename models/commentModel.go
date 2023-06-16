package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid" json:"postID"`
	Post      Post      `gorm:"" json:"post"`
	ProjectID uuid.UUID `gorm:"type:uuid" json:"projectID"`
	Project   Project   `gorm:"" json:"project"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"commentedByID"`
	User      User      `gorm:"" json:"commentedBy"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type UserCommentLike struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"likedByID"`
	User      User      `gorm:"" json:"likedBy"`
	CommentID uuid.UUID `gorm:"type:uuid;not null" json:"commentID"`
	Comment   Comment   `gorm:"" json:"comment"`
	CreatedAt time.Time `json:"likedAt"`
}
