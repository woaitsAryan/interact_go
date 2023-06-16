package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	NotificationType int       `json:"notificationType"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User             User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	PostID           uuid.UUID `gorm:"type:uuid" json:"postID"`
	Post             *Post     `gorm:"constraint:OnDelete:CASCADE" json:"post"`
	ProjectID        uuid.UUID `gorm:"type:uuid" json:"projectID"`
	Project          *Project  `gorm:"constraint:OnDelete:CASCADE" json:"project"`
	OpeningID        uuid.UUID `gorm:"type:uuid" json:"openingID"`
	Opening          *Opening  `gorm:"constraint:OnDelete:CASCADE" json:"opening"`
	CreatedAt        time.Time `json:"createdAt"`
}
