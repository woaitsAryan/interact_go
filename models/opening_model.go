package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Opening struct {
	ID                   uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID            uuid.UUID             `gorm:"type:uuid;not null" json:"projectID"`
	Project              Project               `gorm:"" json:"project"`
	Title                string                `gorm:"type:text;not null" json:"title"`
	Description          string                `gorm:"type:text;not null" json:"description"`
	Tags                 pq.StringArray        `gorm:"type:text[]" json:"tags"`
	Active               bool                  `gorm:"default:true" json:"active"`
	UserID               uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	User                 User                  `gorm:"" json:"user"`
	CreatedAt            time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	NoOfApplications     int16                 `json:"noApplications"`
	Impressions          int                   `gorm:"default:0" json:"noImpressions"`
	Application          []Application         `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"applications"`
	Notifications        []Notification        `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
	LastViews            []LastViewedOpenings  `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
	ProjectHistories     []ProjectHistory      `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
	OpeningBookmarkItems []OpeningBookmarkItem `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
	Messages             []Message             `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
	GroupChatMessages    []GroupChatMessage    `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
	Reports              []Report              `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE" json:"-"`
}
