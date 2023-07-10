package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ChatID    uuid.UUID  `gorm:"type:uuid;not null" json:"chatID"`
	Chat      Chat       `gorm:"" json:"chat"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"userID"`
	User      User       `gorm:"" json:"user"`
	PostID    *uuid.UUID `gorm:"type:uuid" json:"postID"`
	Post      Post       `json:"post"`
	ProjectID *uuid.UUID `gorm:"type:uuid" json:"projectID"`
	Project   Project    `json:"project"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time  `gorm:"default:current_timestamp" json:"createdAt"`
	Read      bool       `gorm:"default:false" json:"read"`
}

type GroupMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ChatID    uuid.UUID `gorm:"type:uuid;not null" json:"chatID"`
	Chat      GroupChat `gorm:"" json:"chat"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"" json:"user"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	Read      bool      `gorm:"default:false" json:"read"`
	ReadBy    []User    `gorm:"many2many:message_read_by;constraint:OnDelete:CASCADE" json:"readBy"`
}

type ProjectChatMessage struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectChatID uuid.UUID   `gorm:"type:uuid;not null" json:"projectChatID"`
	ProjectChat   ProjectChat `gorm:"" json:"projectChat"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null" json:"userID"`
	User          User        `gorm:"" json:"user"`
	Content       string      `gorm:"type:text;not null" json:"content"`
	CreatedAt     time.Time   `gorm:"default:current_timestamp" json:"createdAt"`
	Read          bool        `gorm:"default:false" json:"read"`
	ReadBy        []User      `gorm:"many2many:message_read_by;constraint:OnDelete:CASCADE" json:"readBy"`
}
