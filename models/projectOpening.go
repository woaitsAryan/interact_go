package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Opening struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID        uuid.UUID      `gorm:"type:uuid;not null" json:"projectID"`
	Project          Project        `gorm:"" json:"project"`
	Title            string         `gorm:"type:varchar(255);not null" json:"title"`
	Description      string         `gorm:"type:text;not null" json:"description"`
	Tags             pq.StringArray `gorm:"type:text[]" json:"tags"`
	Active           bool           `gorm:"default:true" json:"active"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"postedByID"`
	User             User           `gorm:"" json:"postedBy"`
	CreatedAt        time.Time      `json:"postedAt"`
	NoOfApplications int            `json:"noOfApplications"`
	Application      []Application  `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"applications,omitempty"`
}
