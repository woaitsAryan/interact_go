package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct { //TODO notifications will be sent after reporting
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ReportType int        `json:"reportType"`
	ReporterID uuid.UUID  `gorm:"type:uuid;not null" json:"reporterID"`
	Reporter   User       `json:"report"`
	UserID     *uuid.UUID `gorm:"type:uuid;" json:"userID"`
	User       User       `json:"user"`
	PostID     *uuid.UUID `gorm:"type:uuid" json:"postID"`
	Post       Post       `json:"post"`
	ProjectID  *uuid.UUID `gorm:"type:uuid" json:"projectID"`
	Project    Project    `json:"project"`
	EventID    *uuid.UUID `gorm:"type:uuid" json:"eventID"`
	Event      Event      `gorm:"" json:"event"`
	OpeningID  *uuid.UUID `gorm:"type:uuid" json:"openingID"`
	Opening    Opening    `json:"opening"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `gorm:"default:current_timestamp" json:"createdAt"`
}
