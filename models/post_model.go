package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID              uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	User                User                  `gorm:"" json:"user"`
	Content             string                `gorm:"type:text;not null" json:"content"`
	CreatedAt           time.Time             `gorm:"default:current_timestamp" json:"postedAt"`
	Images              pq.StringArray        `gorm:"type:text[]" json:"images"`
	Hashes              pq.StringArray        `gorm:"type:text[]" json:"hashes"`
	NoShares            int                   `gorm:"default:0" json:"noShares"`
	NoLikes             int                   `gorm:"default:0" json:"noLikes"`
	NoComments          int                   `gorm:"default:0" json:"noComments"`
	RePostID            *uuid.UUID            `gorm:"type:uuid" json:"rePostID"`
	RePost              *Post                 `gorm:"foreignKey:RePostID" json:"rePost"`
	IsRePost            bool                  `gorm:"default:false" json:"isRePost"`
	NoOfReposts         int                   `gorm:"default:0" json:"noReposts"`
	Tags                pq.StringArray        `gorm:"type:text[]" json:"tags"`
	Impressions         int                   `gorm:"default:0" json:"noImpressions"`
	Edited              bool                  `gorm:"default:false" json:"isEdited"`
	Comments            []Comment             `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments"`
	TaggedUsers         []User                `gorm:"many2many:post_tagged_users" json:"taggedUsers"`
	Notifications       []Notification        `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Messages            []Message             `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	GroupChatMessages   []GroupChatMessage    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Likes               []Like                `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	BookMarkItems       []PostBookmarkItem    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Reports             []Report              `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	OrganizationHistory []OrganizationHistory `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}
