package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID   uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	User     User           `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Content  string         `gorm:"type:text;not null" json:"content"`
	PostedAt time.Time      `json:"postedAt"`
	LikedBy  []*User        `gorm:"many2many:user_post_likes;joinForeignKey:user_id;joinReferences:id" json:"likedBy,omitempty"`
	Images   pq.StringArray `gorm:"type:varchar(255);not null" json:"images"`
	Hashes   pq.StringArray `gorm:"type:text[]" json:"hashes"`
	Caption  string         `gorm:"type:varchar(255);not null" json:"caption"`
	PostedBy uuid.UUID      `gorm:"type:uuid;not null" json:"postedBy"`
	NoShares int            `json:"noShares"`
	NoLikes  int            `json:"noLikes"`
	Tags     pq.StringArray `gorm:"type:text[]" json:"tags"`
	Edited   bool           `gorm:"default:false" json:"edited"`
}
