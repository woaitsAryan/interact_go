package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null" json:"userID"`
	User           User         `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	OrganizationID *uuid.UUID   `gorm:"type:uuid;" json:"organizationID"`
	Organization   Organization `gorm:"constraint:OnDelete:CASCADE" json:"organization"`
	ProjectID      *uuid.UUID   `gorm:"type:uuid;" json:"projectID"`
	Project        Project      `gorm:"constraint:OnDelete:CASCADE" json:"project"`
	Title          string       `gorm:"type:varchar(25);not null" json:"title"`
	Status         int          `gorm:"default:0" json:"status"` //* -1 for reject, 0 for waiting and, 1 for accept
	Read           bool         `gorm:"default:false" json:"isRead"`
	CreatedAt      time.Time    `gorm:"default:current_timestamp" json:"createdAt"`
}
