package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct { //TODO add accepting project invitations field on user model
	//TODO add sender id
	ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID              uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	User                User                  `gorm:"" json:"user"`
	OrganizationID      *uuid.UUID            `gorm:"type:uuid;" json:"organizationID"`
	Organization        Organization          `gorm:"" json:"organization"`
	ProjectID           *uuid.UUID            `gorm:"type:uuid;" json:"projectID"`
	Project             Project               `gorm:"" json:"project"`
	GroupChatID         *uuid.UUID            `gorm:"type:uuid;" json:"chatID"`
	GroupChat           GroupChat             `gorm:"" json:"chat"`
	Title               string                `gorm:"type:varchar(25);not null" json:"title"`
	Status              int                   `gorm:"default:0" json:"status"`     //* -1 for reject, 0 for waiting and, 1 for accept
	Read                bool                  `gorm:"default:false" json:"isRead"` //TODO remove this, not needed
	CreatedAt           time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	OrganizationHistory []OrganizationHistory `gorm:"foreignKey:InvitationID;constraint:OnDelete:CASCADE" json:"-"`
}
