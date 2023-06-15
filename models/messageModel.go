package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ChatID  uuid.UUID `gorm:"type:uuid;not null" json:"chatId"`
	Chat    Chat      `gorm:"constraint:OnDelete:CASCADE" json:"chat"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"senderId"`
	User    User      `gorm:"constraint:OnDelete:CASCADE" json:"sentBy"`
	Content string    `gorm:"type:text;not null" json:"content"`
	SentAt  time.Time `json:"sentAt"`
	Read    bool      `gorm:"default:false" json:"read"`
	ReadBy  []User    `gorm:"many2many:message_read_by;constraint:OnDelete:CASCADE" json:"readBy"`
}
