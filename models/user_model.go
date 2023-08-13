package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	ID                        uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name                      string            `gorm:"varchar(25);not null" json:"name"`
	Username                  string            `gorm:"varchar(10);unique;not null" json:"username"`
	Email                     string            `gorm:"unique;not null" json:"email"`
	Password                  string            `json:"-"`
	ProfilePic                string            `gorm:"default:default.jpg" json:"profilePic"`
	CoverPic                  string            `gorm:"default:default.jpg" json:"coverPic"`
	PhoneNo                   string            `json:"phoneNo"`
	Bio                       string            `json:"bio"`
	Title                     string            `json:"title"`
	Tagline                   string            `json:"tagline"`
	Tags                      pq.StringArray    `gorm:"type:text[]" json:"tags"`
	PasswordResetToken        string            `json:"-"`
	PasswordResetTokenExpires time.Time         `json:"-"`
	Views                     int               `json:"views"`
	NoFollowing               int               `gorm:"default:0" json:"noFollowing"`
	NoFollowers               int               `gorm:"default:0" json:"noFollowers"`
	PasswordChangedAt         time.Time         `gorm:"default:current_timestamp" json:"-"`
	DeactivatedAt             time.Time         `gorm:"default:current_timestamp" json:"-"`
	Admin                     bool              `gorm:"default:false" json:"-"`
	Verified                  bool              `gorm:"default:false" json:"isVerified"`
	LastLoggedIn              time.Time         `gorm:"default:current_timestamp" json:"-"`
	Active                    bool              `gorm:"default:true" json:"-"` //! add a functionality that on delete the acc goes inActive and if the user logs in within 30 days, it goes active again
	CreatedAt                 time.Time         `gorm:"default:current_timestamp" json:"-"`
	Projects                  []Project         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"projects"`
	Posts                     []Post            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts"`
	Memberships               []Membership      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"memberships"`
	Achievements              []Achievement     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"achievements"`
	Applications              []Application     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"applications"`
	PostBookmarks             []PostBookmark    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"postBookmarks"`
	ProjectBookmarks          []ProjectBookmark `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"projectBookmarks"`
	Notifications             []Notification    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"notifications"`
	LastViewed                []LastViewed      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"lastViewed"`
	SendNotifications         []Notification    `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"-"`
	Followers                 []FollowFollower  `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"-"`
	Following                 []FollowFollower  `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE" json:"-"`
	Verification              UserVerification  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (u *User) AfterFind(tx *gorm.DB) error {
	if !u.Active {
		u.Username = "deactived"
		u.Name = "Interact User"
		u.CoverPic = "default.jpg"
		u.ProfilePic = "default.jpg"
		u.Bio = ""
		u.Title = ""
		u.Tagline = ""
		u.Tags = nil
		u.Views = 0
		u.NoFollowers = 0
		u.NoFollowing = 0
		u.Achievements = nil
		u.Memberships = nil
		u.Projects = nil
		u.Posts = nil
	}
	return nil
}

type ProfileView struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	Date   time.Time `gorm:"type:date" json:"date"`
	Count  int       `json:"count"`
}

type FollowFollower struct { //* follower follows followed
	FollowerID uuid.UUID `json:"followerID"`
	Follower   User      `gorm:"foreignKey:FollowerID" json:"follower"`
	FollowedID uuid.UUID `json:"followedID"`
	Followed   User      `gorm:"foreignKey:FollowedID" json:"followed"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"-"`
}

type Achievement struct {
	ID     uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID uuid.UUID      `gorm:"type:uuid;not null" json:"userID"`
	User   User           `json:"user"`
	Title  string         `gorm:"type:text;not null" json:"title"`
	Skills pq.StringArray `gorm:"type:text[];not null" json:"skills"`
}

type LastViewed struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `json:"user"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Project   Project   `json:"project"`
	Timestamp time.Time `gorm:"default:current_timestamp" json:"timestamp"`
}

type UserVerification struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;unique" json:"userID"`
	Code           string    `json:"code"`
	ExpirationTime time.Time `json:"expirationTime"`
}
