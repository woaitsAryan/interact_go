package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Recruitment struct {
	ID             uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OrganizationID uuid.UUID         `gorm:"type:uuid;not null" json:"orgID"`
	Organization   Organization      `gorm:"" json:"organization"`
	Title          string            `gorm:"type:varchar(255);not null" json:"title"`
	Description    string            `gorm:"type:text;not null" json:"description"`
	Tags           pq.StringArray    `gorm:"type:text[]" json:"tags"`
	Tests          []RecruitmentTest `gorm:"foreignKey:RecruitmentID;constraint:OnDelete:CASCADE" json:"openings"`
	CreatedAt      time.Time         `gorm:"default:current_timestamp" json:"createdAt"`
}

type RecruitmentTest struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	RecruitmentID uuid.UUID `gorm:"type:uuid;not null" json:"recruitmentID"`
	CreatedAt     time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}
