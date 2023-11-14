package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Application struct {
	ID            uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OpeningID     uuid.UUID        `gorm:"type:uuid;not null" json:"openingID"`
	Opening       Opening          `gorm:"" json:"opening"`
	ProjectID     uuid.UUID        `gorm:"type:uuid;not null" json:"-"`
	Project       Project          `gorm:"" json:"project"`
	UserID        uuid.UUID        `gorm:"type:uuid;not null" json:"userID"`
	User          User             `gorm:"" json:"user"`
	Email         string           `gorm:"" json:"email"`
	CreatedAt     time.Time        `gorm:"default:current_timestamp" json:"createdAt"`
	Status        int              `json:"status"` //* -1 rejected, 0 submitted, 1 under review, 2 accepted
	Content       string           `gorm:"type:text;not null" json:"content"`
	Resume        string           `gorm:"type:varchar(255)" json:"resume"`
	Links         pq.StringArray   `gorm:"type:text[]" json:"links"`
	Notifications []Notification   `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE" json:"-"`
	History       []ProjectHistory `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE" json:"-"`
}
