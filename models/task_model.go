package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Priority string

const (
	LowPriority    Priority = "low"
	MediumPriority Priority = "medium"
	HighPriority   Priority = "high"
)

type Task struct {
	ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID           *uuid.UUID            `gorm:"" json:"projectID"`
	Project             Project               `gorm:"" json:"project"`
	OrganizationID      *uuid.UUID            `gorm:"" json:"orgID"`
	Organization        Organization          `gorm:"" json:"organization"`
	Deadline            time.Time             `gorm:"default:current_timestamp" json:"deadline"`
	Title               string                `gorm:"type:text;not null" json:"title"`
	Description         string                `gorm:"type:text" json:"description"`
	Tags                pq.StringArray        `gorm:"type:text[]" json:"tags"`
	Users               []User                `gorm:"many2many:task_assigned_users" json:"users"`
	SubTasks            []SubTask             `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"subTasks"`
	Priority            Priority              `gorm:"type:text;default:low" json:"priority"`
	IsCompleted         bool                  `gorm:"default:false" json:"isCompleted"`
	CreatedAt           time.Time             `gorm:"default:current_timestamp;index:idx_created_at,sort:desc" json:"createdAt"`
	OrganizationHistory []OrganizationHistory `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"-"`
}

type SubTask struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	TaskID      uuid.UUID      `gorm:"type:uuid;not null" json:"taskID"`
	Task        Task           `gorm:"" json:"task"`
	Deadline    time.Time      `gorm:"default:current_timestamp;index:idx_deadline,sort:asc" json:"deadline"`
	Title       string         `gorm:"type:text;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Users       []User         `gorm:"many2many:sub_task_assigned_users" json:"users"`
	Priority    Priority       `gorm:"type:text;default:low" json:"priority"`
	IsCompleted bool           `gorm:"default:false" json:"isCompleted"`
}
