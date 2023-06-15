package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatInvitation struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	ChatID    uuid.UUID `gorm:"type:uuid;not null" json:"chatID"`
	Chat      Chat      `gorm:"constraint:OnDelete:CASCADE" json:"chat"`
	Status    int       `gorm:"default:0" json:"status"` //* -1 for reject, 0 for waiting and, 1 for accept
	CreatedAt time.Time `json:"createdAt"`
}

type ProjectInvitation struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Project   Project   `gorm:"constraint:OnDelete:CASCADE" json:"project"`
	Status    int       `gorm:"default:0" json:"status"` //* -1 for reject, 0 for waiting and, 1 for accept
	CreatedAt time.Time `json:"createdAt"`
}
