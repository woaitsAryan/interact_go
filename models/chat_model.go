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
	LastResetByCreatingUser          time.Time `gorm:"default:current_timestamp" json:"-"`
	LastResetByAcceptingUser         time.Time `gorm:"default:current_timestamp" json:"-"`
	BlockedByCreatingUser            bool      `gorm:"default:false" json:"blockedByCreatingUser"`
	BlockedByAcceptingUser           bool      `gorm:"default:false" json:"blockedByAcceptingUser"`
	Messages                         []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID                  uuid.UUID `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage                    *Message  `gorm:"foreignKey:LatestMessageID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Accepted                         bool      `gorm:"default:false" json:"accepted"`
	LastReadMessageByCreatingUserID  uuid.UUID `gorm:"type:uuid" json:"lastReadMessageByCreatingUserID"`
	LastReadMessageByAcceptingUserID uuid.UUID `gorm:"type:uuid" json:"lastReadMessageByAcceptingUserID"`
	LastReadMessageByCreatingUser    *Message  `gorm:"type:uuid" json:"lastReadMessageByCreatingUser"`
	LastReadMessageByAcceptingUser   *Message  `gorm:"type:uuid" json:"lastReadMessageByAcceptingUser"`
}

type GroupChat struct { //TODO store number of members in model to show in invitation
	ID              uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string                `gorm:"type:varchar(50);" json:"title"`
	Description     string                `gorm:"type:text" json:"description"`
	AdminOnly       bool                  `gorm:"default:false" json:"adminOnly"`
	CoverPic        string                `gorm:"type:text; default:default.jpg" json:"coverPic"`
	UserID          uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	User            User                  `gorm:"" json:"user"`
	OrganizationID  *uuid.UUID            `gorm:"type:uuid;" json:"organizationID"`
	Organization    Organization          `gorm:"constraint:OnDelete:CASCADE" json:"organization"`
	ProjectID       *uuid.UUID            `gorm:"type:uuid;" json:"projectID"`
	Project         Project               `gorm:"constraint:OnDelete:CASCADE" json:"project"`
	CreatedAt       time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	Memberships     []GroupChatMembership `gorm:"constraint:OnDelete:CASCADE" json:"memberships"`
	Messages        []GroupChatMessage    `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID uuid.UUID             `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage   *GroupChatMessage     `gorm:"foreignKey:LatestMessageID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Invitations     []Invitation          `gorm:"foreignKey:GroupChatID;constraint:OnDelete:CASCADE" json:"invitations"`
}

type GroupChatRole string

const (
	ChatMember GroupChatRole = "Member"
	ChatAdmin  GroupChatRole = "Admin"
)

type GroupChatMembership struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID      uuid.UUID     `gorm:"type:uuid;not null" json:"userID"`
	User        User          `gorm:"" json:"user"`
	Role        GroupChatRole `gorm:"type:text" json:"role"`
	GroupChatID uuid.UUID     `gorm:"type:uuid;not null" json:"chatID"`
	GroupChat   GroupChat     `gorm:"" json:"chat"`
	CreatedAt   time.Time     `gorm:"default:current_timestamp" json:"createdAt"`
}
