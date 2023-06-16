package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Application struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OpeningID uuid.UUID      `gorm:"type:uuid;not null" json:"openingID"`
	Opening   Opening        `gorm:"" json:"opening"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"applicantID"`
	User      User           `gorm:"" json:"applicant"`
	CreatedAt time.Time      `json:"appliedAt"`
	Status    int            `json:"status"` //* -1 rejected, 0 submitted, 1 under review, 2 accepted
	Content   string         `gorm:"type:text;not null" json:"content"`
	Resume    string         `gorm:"type:varchar(255)" json:"resume"`
	Links     pq.StringArray `gorm:"type:text[]" json:"links"`
}
