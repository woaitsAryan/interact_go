package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Project struct {
	ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title               string                `gorm:"type:text;not null" json:"title"` //TODO Validation error handling for no of chars
	Slug                string                `gorm:"type:text;not null" json:"slug"`
	Tagline             string                `gorm:"type:text;not null" json:"tagline"`
	CoverPic            string                `gorm:"type:text; default:default.jpg" json:"coverPic"`
	Hash                string                `gorm:"type:text" json:"hash"`
	Description         string                `gorm:"type:text;not null" json:"description"`
	Page                string                `gorm:"type:text" json:"page"`
	UserID              uuid.UUID             `gorm:"type:uuid;not null" json:"userID"`
	User                User                  `gorm:"" json:"user"`
	CreatedAt           time.Time             `gorm:"default:current_timestamp" json:"createdAt"`
	Tags                pq.StringArray        `gorm:"type:text[]" json:"tags"`
	NoLikes             int                   `gorm:"default:0" json:"noLikes"`
	NoShares            int                   `gorm:"default:0" json:"noShares"`
	NoComments          int                   `gorm:"default:0" json:"noComments"`
	TotalNoViews        int                   `gorm:"default:0" json:"totalNoViews"`
	Category            string                `gorm:"type:text;not null" json:"category"`
	IsPrivate           bool                  `gorm:"default:false" json:"isPrivate"`
	TRatio              int                   `json:"-"`
	Views               int                   `json:"views"`
	NumberOfMembers     int                   `gorm:"default:1" json:"noMembers"`
	Impressions         int                   `gorm:"default:1" json:"impressions"`
	Links               pq.StringArray        `gorm:"type:text[]" json:"links"`
	PrivateLinks        pq.StringArray        `gorm:"type:text[]" json:"-"`
	Comments            []Comment             `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"comments"`
	Openings            []Opening             `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"openings"`
	Chats               []GroupChat           `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"chats"`
	Invitations         []Invitation          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"invitations"`
	Memberships         []Membership          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"memberships"`
	Tasks               []Task                `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"tasks"`
	History             []ProjectHistory      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Notifications       []Notification        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	LastViews           []LastViewedProjects  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Messages            []Message             `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	GroupChatMessages   []GroupChatMessage    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	BookMarkItems       []ProjectBookmarkItem `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	ProjectViews        []ProjectView         `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Likes               []Like                `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Applications        []Application         `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Reports             []Report              `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	OrganizationHistory []OrganizationHistory `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

type ProjectView struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Date      time.Time `json:"date"`
	Count     int       `json:"count"`
}

/*
history type:
*-1 - User created this project
*0 - User sent invitation to user
*1 - User joined this project
*2 - User edited project details
*3 - User created an opening
*4 - User edited opening details
*5 - User deleted opening
*6 - User accepted application of user
*7 - User rejected application of user
*8 - User created a new group chat
*9 - User created a new task
*10 - User left the project
*11 - User removed user from the project
*/

type ProjectHistory struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID     uuid.UUID   `gorm:"type:uuid;not null" json:"projectID"`
	HistoryType   int         `json:"historyType"`
	SenderID      uuid.UUID   `gorm:"type:uuid;not null" json:"senderID"`
	Sender        User        `json:"sender"`
	UserID        *uuid.UUID  `gorm:"type:uuid" json:"userID"`
	User          User        `json:"user"`
	OpeningID     *uuid.UUID  `gorm:"type:uuid" json:"openingID"`
	Opening       Opening     `json:"opening"`
	ApplicationID *uuid.UUID  `gorm:"type:uuid" json:"applicationID"`
	Application   Application `json:"application"`
	InvitationID  *uuid.UUID  `gorm:"type:uuid" json:"invitationID"`
	Invitation    Invitation  `json:"invitation"`
	TaskID        *uuid.UUID  `gorm:"type:uuid" json:"taskID"`
	Task          Task        `json:"task"`
	DeletedText   string      `gorm:"type:text" json:"deletedMessage"`
	CreatedAt     time.Time   `gorm:"default:current_timestamp;index:idx_created_at,sort:desc" json:"createdAt"`
}
