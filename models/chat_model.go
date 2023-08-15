package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CreatingUserID  uuid.UUID `gorm:"type:uuid;not null" json:"createdByID"`
	CreatingUser    User      `gorm:"" json:"createdBy"`
	AcceptingUserID uuid.UUID `gorm:"type:uuid;not null" json:"acceptedByID"`
	AcceptingUser   User      `gorm:"" json:"acceptedBy"`
	CreatedAt       time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	Messages        []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID uuid.UUID `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage   *Message  `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Accepted        bool      `gorm:"default:false" json:"accepted"`
}

type GroupChat struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string         `gorm:"type:varchar(50);" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	CreatingUserID  uuid.UUID      `gorm:"type:uuid;not null" json:"createdByID"`
	CreatingUser    User           `gorm:"" json:"createdBy"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
	Members         []User         `gorm:"many2many:chat_members;constraint:OnDelete:CASCADE" json:"members"`
	Messages        []GroupMessage `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID uuid.UUID      `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage   *GroupMessage  `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Invitations     []Invitation   `gorm:"foreignKey:GroupChatID;constraint:OnDelete:CASCADE" json:"invitations"`
}

type ProjectChat struct {
	ID              uuid.UUID               `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string                  `gorm:"type:varchar(50);" json:"title"`
	Description     string                  `gorm:"type:text" json:"description"`
	CreatingUserID  uuid.UUID               `gorm:"type:uuid;not null" json:"createdByID"`
	CreatingUser    User                    `gorm:"" json:"createdBy"`
	ProjectID       uuid.UUID               `gorm:"type:uuid;not null" json:"projectID"`
	Project         Project                 `gorm:"" json:"project"`
	CreatedAt       time.Time               `gorm:"default:current_timestamp" json:"createdAt"`
	Memberships     []ProjectChatMembership `gorm:"constraint:OnDelete:CASCADE" json:"memberships"`
	LatestMessageID uuid.UUID               `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage   *ProjectChatMessage     `gorm:"foreignKey:ProjectChatID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Messages        []ProjectChatMessage    `gorm:"foreignKey:ProjectChatID;constraint:OnDelete:CASCADE" json:"messages"`
}

type ProjectChatMembership struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null" json:"userID"`
	User          User        `gorm:"" json:"user"`
	ProjectID     uuid.UUID   `gorm:"type:uuid;not null" json:"projectID"`
	Project       Project     `gorm:"" json:"project"`
	ProjectChatID uuid.UUID   `gorm:"type:uuid;not null" json:"projectChatID"`
	ProjectChat   ProjectChat `gorm:"" json:"projectChat"`
	CreatedAt     time.Time   `gorm:"default:current_timestamp" json:"createdAt"`
}
