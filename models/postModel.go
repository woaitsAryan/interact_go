package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"userID"`
	User      User           `gorm:"" json:"user"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"postedAt"`
	LikedBy   []User         `gorm:"many2many:user_post_likes;joinForeignKey:user_id;joinReferences:id;constraint:OnDelete:CASCADE" json:"likedBy,omitempty"`
	Images    pq.StringArray `gorm:"type:text[]" json:"images"`
	Hashes    pq.StringArray `gorm:"type:text[]" json:"hashes"`
	NoShares  int            `json:"noShares"`
	NoLikes   int            `json:"noLikes"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	Edited    bool           `gorm:"default:false" json:"edited"`
	Comments  []PostComment  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}
