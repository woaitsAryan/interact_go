package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"likedByID"`
	PostID    *uuid.UUID `gorm:"type:uuid" json:"postID"`
	ProjectID *uuid.UUID `gorm:"type:uuid" json:"projectID"`
	EventID   *uuid.UUID `gorm:"type:uuid" json:"eventID"`
	CommentID *uuid.UUID `gorm:"type:uuid" json:"commentID"`
	CreatedAt time.Time  `gorm:"default:current_timestamp" json:"likedAt"`
}
