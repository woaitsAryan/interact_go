package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ChatID   uuid.UUID `gorm:"type:uuid;not null" json:"chatId"`
	Chat     Chat      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	SenderID uuid.UUID `gorm:"type:uuid;not null" json:"senderId"`
	Sender   User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Text     string    `gorm:"type:text;not null" json:"text"`
	SentAt   time.Time `json:"sentAt"`
	IsRead   bool      `gorm:"default:false" json:"isRead"`
	ReadBy   []User    `gorm:"many2many:message_read_by;constraint:OnDelete:CASCADE" json:"readBy"`
}
