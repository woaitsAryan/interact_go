package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Project struct {
	ID                     uuid.UUID               `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title                  string                  `gorm:"type:varchar(255);not null" json:"title"` //TODO Validation error handling for no of chars
	Tagline                string                  `gorm:"type:varchar(255);not null" json:"tagline"`
	CoverPic               string                  `gorm:"type:varchar(255); default:default.jpg" json:"coverPic"`
	Hash                   string                  `gorm:"type:varchar(255)" json:"hash"`
	Description            string                  `gorm:"type:text;not null" json:"description"`
	Page                   string                  `gorm:"type:text" json:"page"`
	UserID                 uuid.UUID               `gorm:"type:uuid;not null" json:"userID"`
	User                   User                    `gorm:"" json:"user"`
	CreatedAt              time.Time               `gorm:"default:current_timestamp" json:"createdAt"`
	Tags                   pq.StringArray          `gorm:"type:text[]" json:"tags"`
	NoLikes                int                     `json:"noLikes"`
	NoShares               int                     `json:"noShares"`
	NoComments             int                     `gorm:"default:0" json:"noComments"`
	Category               string                  `gorm:"type:varchar(255);not null" json:"category"`
	IsPrivate              bool                    `gorm:"default:false" json:"isPrivate"`
	TRatio                 int                     `json:"-"`
	Views                  int                     `json:"views"`
	Links                  pq.StringArray          `gorm:"type:text[]" json:"links"`
	PrivateLinks           pq.StringArray          `gorm:"type:text[]" json:"privateLinks"`
	Comments               []Comment               `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"comments"`
	Openings               []Opening               `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"openings"`
	Chats                  []ProjectChat           `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"chats"`
	Invitations            []Invitation            `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"invitations"`
	Memberships            []Membership            `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"memberships"`
	Notifications          []Notification          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"notifications"`
	LastViews              []LastViewed            `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Messages               []Message               `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	BookMarkItems          []ProjectBookmarkItem   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	ProjectChatMemberships []ProjectChatMembership `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	ProjectViews           []ProjectView           `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Likes                  []UserProjectLike       `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Applications           []Application           `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	CommentLikes           []UserCommentLike       `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	ChatMessages           []ProjectChatMessage    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

type ProjectView struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Date      time.Time `json:"date"`
	Count     int       `json:"count"`
}

type UserProjectLike struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"likedByID"`
	User      User      `gorm:"" json:"likedBy"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Project   Project   `gorm:"" json:"project"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"likedAt"`
}
