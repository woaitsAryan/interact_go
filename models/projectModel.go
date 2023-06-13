package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Project struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	Tagline      string         `gorm:"type:varchar(255);not null" json:"tagline"`
	CoverPic     string         `gorm:"type:varchar(255)" json:"coverPic"`
	Hashes       pq.StringArray `gorm:"type:text[]" json:"hashes"`
	Description  string         `gorm:"type:text;not null" json:"description"`
	Page         string         `gorm:"type:text" json:"page"`
	PostedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"postedBy"`
	User         User           `gorm:"foreignKey:PostedBy;references:ID" json:"-"`
	PostedAt     time.Time      `json:"postedAt"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	NoLikes      int            `json:"noLikes"`
	NoShares     int            `json:"noShares"`
	Category     string         `gorm:"type:varchar(255);not null" json:"category"`
	IsPrivate    bool           `gorm:"default:false" json:"isPrivate"`
	TRatio       int            `json:"-"`
	Links        pq.StringArray `gorm:"type:text[]" json:"links"`
	PrivateLinks pq.StringArray `gorm:"type:text[]" json:"privateLinks"`
	LikedBy      []*User        `gorm:"many2many:user_project_likes;joinForeignKey:user_id;joinReferences:id" json:"likedBy,omitempty"`
}

type ProjectView struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectId"`
	Date      time.Time `json:"date"`
	Count     int       `json:"count"`
}
