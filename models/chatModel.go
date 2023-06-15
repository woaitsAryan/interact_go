package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title       string    `gorm:"type:varchar(50);" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"createdByID"`
	User        User      `gorm:"constraint:OnDelete:CASCADE" json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	Members     []User    `gorm:"many2many:chat_members;constraint:OnDelete:CASCADE" json:"members"`
	Group       bool      `gorm:"default:false" json:"group"`
}

type ProjectChat struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title       string    `gorm:"type:varchar(50);" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"createdByID"`
	User        User      `gorm:"constraint:OnDelete:CASCADE" json:"createdBy"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	Project     Project   `gorm:"constraint:OnDelete:CASCADE" json:"project"`
	CreatedAt   time.Time `json:"createdAt"`
	Members     []User    `gorm:"many2many:project_chat_members;constraint:OnDelete:CASCADE" json:"members"`
}
