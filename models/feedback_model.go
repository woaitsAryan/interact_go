package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Type      int8      `json:"type"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}
