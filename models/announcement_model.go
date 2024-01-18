package models

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	OrganizationID      uuid.UUID             `gorm:"type:uuid;not null" json:"organizationID"`
	Organization        Organization          `gorm:"" json:"organization"`
	Title               string                `gorm:"" json:"title"`
	Content             string                `gorm:"not null" json:"content"`
	IsEdited            bool                  `gorm:"default:false" json:"isEdited"`
	IsOpen              bool                  `gorm:"default:false" json:"isOpen"`
	CreatedAt           time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	NoShares            int                   `gorm:"default:0" json:"noShares"`
	NoLikes             int                   `gorm:"default:0" json:"noLikes"`
	NoComments          int                   `gorm:"default:0" json:"noComments"`
	Comments            []Comment             `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	TaggedUsers         []User                `gorm:"many2many:announcement_tagged_users" json:"taggedUsers"`
	Notifications       []Notification        `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	Messages            []Message             `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	GroupChatMessages   []GroupChatMessage    `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	Likes               []Like                `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	Reports             []Report              `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	OrganizationHistory []OrganizationHistory `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
}
