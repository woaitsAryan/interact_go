package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Project struct {
	ID           uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title        string           `gorm:"type:varchar(255);not null" json:"title"`
	Tagline      string           `gorm:"type:varchar(255);not null" json:"tagline"`
	CoverPic     string           `gorm:"type:varchar(255)" json:"coverPic"`
	Hash         string           `gorm:"type:varchar(255)" json:"hash"`
	Description  string           `gorm:"type:text;not null" json:"description"`
	Page         string           `gorm:"type:text" json:"page"`
	UserID       uuid.UUID        `gorm:"type:uuid;not null" json:"userID"`
	User         User             `gorm:"" json:"user"`
	CreatedAt    time.Time        `json:"postedAt"`
	Tags         pq.StringArray   `gorm:"type:text[]" json:"tags"`
	NoLikes      int              `json:"noLikes"`
	NoShares     int              `json:"noShares"`
	Category     string           `gorm:"type:varchar(255);not null" json:"category"`
	IsPrivate    bool             `gorm:"default:false" json:"isPrivate"`
	TRatio       int              `json:"-"`
	Links        pq.StringArray   `gorm:"type:text[]" json:"links"`
	PrivateLinks pq.StringArray   `gorm:"type:text[]" json:"privateLinks"`
	LikedBy      []User           `gorm:"many2many:user_project_likes;joinForeignKey:user_id;joinReferences:id;constraint:OnDelete:CASCADE" json:"likedBy,omitempty"`
	Comments     []ProjectComment `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Openings     []Opening        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"openings,omitempty"`
	Chats        []ProjectChat    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"chats,omitempty"`
	Invitations  []ChatInvitation `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"invitations"`
}

type ProjectView struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectId"`
	Date      time.Time `json:"date"`
	Count     int       `json:"count"`
}
