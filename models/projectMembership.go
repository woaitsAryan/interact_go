package models

import (
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null" json:"projectId"`
	Project    Project   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User       User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Role       string    `gorm:"type:varchar(25);not null" json:"role"`
	Active     bool      `gorm:"default:true" json:"active"`
	DateJoined time.Time
}
