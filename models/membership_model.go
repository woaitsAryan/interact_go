package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectRole string

const (
	ProjectMember  ProjectRole = "Member"
	ProjectEditor  ProjectRole = "Editor"
	ProjectManager ProjectRole = "Manager"
)

type Membership struct {
	ID               uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID        uuid.UUID        `gorm:"type:uuid;not null" json:"projectID"`
	Project          Project          `gorm:"" json:"project"`
	UserID           uuid.UUID        `gorm:"type:uuid;not null" json:"userID"`
	User             User             `gorm:"" json:"user"`
	Role             ProjectRole      `gorm:"type:text;not null" json:"role"`
	Title            string           `gorm:"type:varchar(25);not null" json:"title"`
	Active           bool             `gorm:"default:true" json:"active"`
	CreatedAt        time.Time        `gorm:"default:current_timestamp" json:"createdAt"`
	ProjectHistories []ProjectHistory `gorm:"foreignKey:MembershipID;constraint:OnDelete:CASCADE" json:"-"`
}
