package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ChatID        uuid.UUID   `gorm:"type:uuid" json:"chatID"`
	Chat          Chat        `gorm:"" json:"chat"`
	ProjectChatID uuid.UUID   `gorm:"type:uuid" json:"projectChatID"`
	ProjectChat   ProjectChat `gorm:"" json:"projectChat"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null" json:"senderID"`
	User          User        `gorm:"" json:"sentBy"`
	Content       string      `gorm:"type:text;not null" json:"content"`
	CreatedAt     time.Time   `json:"sentAt"`
	Read          bool        `gorm:"default:false" json:"read"`
	ReadBy        []User      `gorm:"many2many:message_read_by;constraint:OnDelete:CASCADE" json:"readBy"`
}
