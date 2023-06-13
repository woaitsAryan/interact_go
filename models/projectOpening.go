package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Opening struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null" json:"projectId"`
	Project        Project        `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Description    string         `gorm:"type:text;not null" json:"description"`
	NoPositions    int            `gorm:"not null" json:"noPositions"`
	SkillsRequired pq.StringArray `gorm:"type:text[]" json:"skillsRequired"`
	IsActive       bool           `gorm:"default:true" json:"isActive"`
	PostedBy       uuid.UUID      `gorm:"type:uuid;not null" json:"postedBy"`
	User           User           `gorm:"constraint:OnDelete:CASCADE;foreignKey:PostedBy" json:"-"`
	PostedAt       time.Time      `json:"postedAt"`
}
