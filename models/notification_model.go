package models

import (
	"time"

	"github.com/google/uuid"
)

/*
notification type:
*-1 - Welcome to Interact
*0 - User started following you
*1 - User liked your post
*2 - User commented on your post
*3 - User liked your project
*4 - User commented on your project
*5 - User applied for your project opening
*6 - You got selected for the opening
*7 - You got rejected for the opening
*8 - You were removed from the project //TODO have to implement this
*9 - Chat request
*10 - Accepted Project Invitation //TODO add more invitation acceptance notifications, and add notification for you have been invited
*11 - User assigned you a task in project
*12 - User liked your event
*13 - User commented on your event
*14 - Your post got x impressions
*15 - Your project got x impressions
*16 - Your event got x impressions
*17 - Your announcement got x impressions
*18 - User liked your announcement
*19 - User commented on your announcement
*20 - User applied for your organization's opening
*21 - Tagged in a Post
*22 - Tagged in an Announcement
*/

type Notification struct {
	ID               uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	NotificationType int16        `json:"notificationType"`
	UserID           uuid.UUID    `gorm:"type:uuid;not null" json:"userID"`
	User             User         `json:"user"`
	SenderID         uuid.UUID    `gorm:"type:uuid;not null" json:"senderID"`
	Sender           User         `json:"sender"`
	PostID           *uuid.UUID   `gorm:"type:uuid" json:"postID"`
	Post             Post         `json:"post"`
	ProjectID        *uuid.UUID   `gorm:"type:uuid" json:"projectID"`
	Project          Project      `json:"project"`
	EventID          *uuid.UUID   `gorm:"type:uuid" json:"eventID"`
	Event            Event        `json:"event"`
	AnnouncementID   *uuid.UUID   `gorm:"type:uuid" json:"announcementID"`
	Announcement     Announcement `json:"announcement"`
	OpeningID        *uuid.UUID   `gorm:"type:uuid" json:"openingID"`
	Opening          Opening      `json:"opening"`
	ApplicationID    *uuid.UUID   `gorm:"type:uuid" json:"applicationID"`
	Application      Application  `json:"application"`
	ImpressionCount  int          `gorm:"default:0" json:"impressionCount"`
	Read             bool         `gorm:"default:false" json:"isRead"`
	CreatedAt        time.Time    `gorm:"default:current_timestamp" json:"createdAt"`
}
