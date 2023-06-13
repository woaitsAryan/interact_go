package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	CreatedByID uuid.UUID `gorm:"type:uuid;not null" json:"createdBy"`
	CreatedBy   User      `gorm:"constraint:OnDelete:CASCADE;foreignKey:CreatedByID" json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	Members     []User    `gorm:"many2many:chat_members;constraint:OnDelete:CASCADE" json:"groupMembers"`
}
