package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Application struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OpeningID         uuid.UUID      `gorm:"type:uuid;not null" json:"openingId"`
	Opening           Opening        `gorm:"" json:"opening"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null" json:"applicantId"`
	User              User           `gorm:"" json:"applicant"`
	CreatedAt         time.Time      `json:"appliedAt"`
	ApplicationStatus int            `json:"applicationStatus"`
	Content           string         `gorm:"type:text;not null" json:"content"`
	Resume            string         `gorm:"type:varchar(255)" json:"resume"`
	Skills            pq.StringArray `gorm:"type:text[]" json:"skills"`
	Links             pq.StringArray `gorm:"type:text[]" json:"links"`
}
