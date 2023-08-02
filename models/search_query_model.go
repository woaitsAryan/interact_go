package models

import (
	"time"

	"github.com/google/uuid"
)

type SearchQuery struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Query     string    `gorm:"index"`
	Timestamp time.Time `gorm:"default:current_timestamp" json:"-"`
}
