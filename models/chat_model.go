package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID                               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CreatingUserID                   uuid.UUID `gorm:"type:uuid;not null" json:"createdByID"`
	CreatingUser                     User      `gorm:"" json:"createdBy"`
	AcceptingUserID                  uuid.UUID `gorm:"type:uuid;not null" json:"acceptedByID"`
	AcceptingUser                    User      `gorm:"" json:"acceptedBy"`
	CreatedAt                        time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	Messages                         []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID                  uuid.UUID `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage                    *Message  `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Accepted                         bool      `gorm:"default:false" json:"accepted"`
	LastReadMessageByCreatingUserID  uuid.UUID `gorm:"type:uuid" json:"lastReadMessageByCreatingUserID"`
	LastReadMessageByAcceptingUserID uuid.UUID `gorm:"type:uuid" json:"lastReadMessageByAcceptingUserID"`
}

type GroupChat struct {
	ID              uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string                `gorm:"type:varchar(50);" json:"title"`
	Description     string                `gorm:"type:text" json:"description"`
	UserID          uuid.UUID             `gorm:"type:uuid;not null" json:"createdByID"`
	User            User                  `gorm:"" json:"createdBy"`
	OrganizationID  *uuid.UUID            `gorm:"type:uuid;" json:"organizationID"`
	Organization    Organization          `gorm:"constraint:OnDelete:CASCADE" json:"organization"`
	ProjectID       *uuid.UUID            `gorm:"type:uuid;" json:"projectID"`
	Project         Project               `gorm:"constraint:OnDelete:CASCADE" json:"project"`
	CreatedAt       time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	Memberships     []GroupChatMembership `gorm:"constraint:OnDelete:CASCADE" json:"memberships"`
	Messages        []GroupChatMessage    `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID uuid.UUID             `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage   *GroupChatMessage     `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Invitations     []Invitation          `gorm:"foreignKey:GroupChatID;constraint:OnDelete:CASCADE" json:"invitations"`
}
type GroupChatMembership struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User        User      `gorm:"" json:"user"`
	GroupChatID uuid.UUID `gorm:"type:uuid;not null" json:"chatID"`
	GroupChat   GroupChat `gorm:"" json:"chat"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}
