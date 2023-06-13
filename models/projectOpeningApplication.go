package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Application struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	OpeningID         uuid.UUID      `gorm:"type:uuid;not null" json:"openingId"`
	Opening           Opening        `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ApplicantID       uuid.UUID      `gorm:"type:uuid;not null" json:"applicantId"`
	Applicant         User           `gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicantID" json:"-"`
	AppliedAt         time.Time      `json:"appliedAt"`
	ApplicationStatus int            `json:"applicationStatus"`
	Resume            string         `gorm:"type:varchar(255);not null" json:"resume"`
	CoverLetter       string         `gorm:"type:text;not null" json:"coverLetter"`
	Skills            pq.StringArray `gorm:"type:text[]" json:"skills"`
	Experience        string         `json:"experience"`
	Education         string         `json:"education"`
	Portfolio         string         `json:"portfolio"`
	ProjectLinks      pq.StringArray `gorm:"type:text[]" json:"projectLinks"`
}
