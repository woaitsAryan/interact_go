package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"postedByID"`
	User      User      `gorm:"" json:"postedBy"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"postedAt"`
	// LikedBy   []User         `gorm:"many2many:user_post_likes;joinForeignKey:user_id;joinReferences:id;constraint:OnDelete:CASCADE" json:"likedBy,omitempty"`
	Images   []string  `gorm:"type:text[]" json:"images"`
	Hashes   []string  `gorm:"type:text[]" json:"hashes"`
	NoShares int       `json:"noShares"`
	NoLikes  int       `json:"noLikes"`
	Tags     []string  `gorm:"type:text[]" json:"tags"`
	Edited   bool      `gorm:"default:false" json:"edited"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

type UserPostLike struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"likedByID"`
	User      User      `gorm:"" json:"likedBy"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"postID"`
	Post      Post      `gorm:"" json:"post"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"likedAt"`
}
