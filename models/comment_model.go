package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID               uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PostID           *uuid.UUID   `gorm:"type:uuid" json:"postID"`
	Post             Post         `gorm:"" json:"post"`
	ProjectID        *uuid.UUID   `gorm:"type:uuid" json:"projectID"`
	Project          Project      `gorm:"" json:"project"`
	EventID          *uuid.UUID   `gorm:"type:uuid" json:"eventID"`
	Event            Event        `gorm:"" json:"event"`
	AnnouncementID   *uuid.UUID   `gorm:"type:uuid" json:"announcementID"`
	Announcement     Announcement `gorm:"" json:"announcement"`
	ParentCommentID  *uuid.UUID   `gorm:"type:uuid;" json:"parentCommentID"`
	IsRepliedComment bool         `gorm:"default:false" json:"isRepliedComment"`
	UserID           uuid.UUID    `gorm:"type:uuid;not null" json:"userID"`
	User             User         `gorm:"" json:"user"`
	Content          string       `gorm:"type:text;not null" json:"content"`
	NoLikes          int          `json:"noLikes"`
	NoReplies        int          `json:"noReplies"`
	Edited           bool         `gorm:"default:false" json:"edited"`
	IsFlagged        bool         `gorm:"default:false" json:"-"`
	CreatedAt        time.Time    `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt        time.Time    `gorm:"default:current_timestamp" json:"updatedAt"`
	Likes            []Like       `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"-"`
}
