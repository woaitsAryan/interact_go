package models

import (
	"time"

	"github.com/google/uuid"
)

/*
	Options Model

VotedBy array is an array of user IDs who have voted for this option
*/
type Option struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PollID    uuid.UUID `gorm:"type:uuid;not null" json:"pollID"`
	Content   string    `gorm:"not null" json:"content"`
	Votes     int       `gorm:"type:int;default:0" json:"noVotes"`
	VotedBy   []User    `gorm:"many2many:voted_by;constraint:OnDelete:CASCADE" json:"votedBy"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}

/*
	Poll Model

Has a one to many relationship with options
Has a multi answer option with isMultiAnswer field
*/
type Poll struct {
	ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OrganizationID      uuid.UUID             `gorm:"type:uuid;not null" json:"organizationID"`
	Organization        Organization          `gorm:"" json:"organization"`
	Title               string                `gorm:"" json:"title"`
	Content             string                `gorm:"not null" json:"content"`
	Options             []Option              `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE" json:"options"`
	IsMultiAnswer       bool                  `gorm:"default:false" json:"isMultiAnswer"`
	IsEdited            bool                  `gorm:"default:false" json:"isEdited"`
	IsOpen              bool                  `gorm:"default:false" json:"isOpen"`
	TotalVotes          int                   `gorm:"default:0" json:"totalVotes"`
	CreatedAt           time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	OrganizationHistory []OrganizationHistory `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE" json:"-"`
}
