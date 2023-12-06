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
	Invitations       []Invitation             `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"invitations"`
	History           []OrganizationHistory    `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"history"`
	CreatedAt         time.Time                `gorm:"default:current_timestamp" json:"createdAt"`
}

type OrganizationRole string

const (
	Member  OrganizationRole = "Member"
	Senior  OrganizationRole = "Senior"
	Manager OrganizationRole = "Manager"
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
	Title          string           `gorm:"type:varchar(25);not null" json:"title"`
	Role           OrganizationRole `gorm:"type:text" json:"role"`
	CreatedAt      time.Time        `gorm:"default:current_timestamp" json:"createdAt"`
}

type OrganizationHistory struct {
	ID             uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OrganizationID uuid.UUID   `gorm:"type:uuid;not null" json:"orgID"`
	HistoryType    int         `json:"historyType"`
	UserID         uuid.UUID   `gorm:"type:uuid;not null" json:"userID"`
	User           User        `json:"user"`
	PostID         *uuid.UUID  `gorm:"type:uuid" json:"postID"`
	Post           Post        `json:"post"`
	ProjectID      *uuid.UUID  `gorm:"type:uuid" json:"projectID"`
	Project        Project     `json:"project"`
	OpeningID      *uuid.UUID  `gorm:"type:uuid" json:"openingID"`
	Opening        Opening     `json:"opening"`
	ApplicationID  *uuid.UUID  `gorm:"type:uuid" json:"applicationID"`
	Application    Application `json:"application"`
	CreatedAt      time.Time   `gorm:"default:current_timestamp" json:"createdAt"`
}
