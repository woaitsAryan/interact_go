package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID                uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID            uuid.UUID                `gorm:"type:uuid;not null" json:"userID"` //user model who is given the organization status
	User              User                     `gorm:"" json:"user"`
	OrganizationTitle string                   `gorm:"unique" json:"title"`
	Memberships       []OrganizationMembership `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"memberships"`
	CreatedAt         time.Time                `gorm:"default:current_timestamp" json:"createdAt"`
}

type OrganizationRole string

const (
	Member  OrganizationRole = "Member"
	Manager OrganizationRole = "Manager"
	Owner   OrganizationRole = "Owner"
)

//* Member can only view.
//* Manager can CRUD on Projects and Posts.
//* Owner can CRUD on members can change their roles and can update the Organization Details.
//* Organization Account Logger can do all this, including delete the organization, transfer of ownership, and organization chats.

type OrganizationMembership struct {
	ID             uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OrganizationID uuid.UUID        `gorm:"type:uuid;not null" json:"organizationID"`
	Organization   Organization     `gorm:"" json:"organization"`
	UserID         uuid.UUID        `gorm:"type:uuid;not null" json:"userID"`
	User           User             `gorm:"" json:"user"`
	Role           OrganizationRole `gorm:"type:text" json:"role"`
	CreatedAt      time.Time        `gorm:"default:current_timestamp" json:"createdAt"`
}
