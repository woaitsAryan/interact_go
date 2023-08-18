package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostID    *uuid.UUID        `gorm:"type:uuid" json:"postID"`
	Post      Post              `gorm:"" json:"post"`
	ProjectID *uuid.UUID        `gorm:"type:uuid" json:"projectID"`
	Project   Project           `gorm:"" json:"project"`
	UserID    uuid.UUID         `gorm:"type:uuid;not null" json:"userID"`
	User      User              `gorm:"" json:"user"`
	Content   string            `gorm:"type:text;not null" json:"content"`
	NoLikes   int               `json:"noLikes"`
	Edited    bool              `gorm:"default:false" json:"edited"`
	CreatedAt time.Time         `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time         `gorm:"default:current_timestamp" json:"updatedAt"`
	Likes     []UserCommentLike `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserCommentLike struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"userID"`
	User      User       `gorm:"" json:"user"`
	CommentID uuid.UUID  `gorm:"type:uuid;not null" json:"commentID"`
	Comment   Comment    `gorm:"" json:"comment"`
	PostID    *uuid.UUID `gorm:"type:uuid" json:"-"`
	ProjectID *uuid.UUID `gorm:"type:uuid" json:"-"`
	CreatedAt time.Time  `gorm:"default:current_timestamp" json:"createdAt"`
}
