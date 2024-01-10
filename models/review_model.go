package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
	Organizational Review Model

Models a review for an organization.
Has an anonymous field to review anonymously.
Has a relevance field to compute relevance of the review and sort by it.
*/
type Review struct {
	ID                uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID            uuid.UUID    `gorm:"type:uuid;not null" json:"userID"`
	User              User         `gorm:"" json:"user"`
	OrganizationID    uuid.UUID    `gorm:"type:uuid;not null" json:"organizationID"`
	Organization      Organization `gorm:"" json:"-"`
	Content           string       `gorm:"type:text;not null" json:"content"`
	Relevance         int8         `gorm:"not null;default:10" json:"-"`
	Rating            int8         `gorm:"not null;default:0" json:"rating"`
	NumberOfUpVotes   int          `gorm:"not null;default:0" json:"noUpVotes"`
	NumberOfDownVotes int          `gorm:"not null;default:0" json:"noDownVotes"`
	Anonymous         bool         `gorm:"not null;default:false" json:"isAnonymous"`
	CreatedAt         time.Time    `gorm:"default:current_timestamp" json:"createdAt"`
}

func (r *Review) AfterFind(tx *gorm.DB) error {
	if r.Anonymous {
		r.UserID = uuid.Nil
		r.User = User{}
	}
	return nil
}
