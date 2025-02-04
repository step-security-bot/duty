package models

import (
	"time"

	"github.com/google/uuid"
)

// Artifact represents the artifacts table
type Artifact struct {
	ID            uuid.UUID  `json:"id" gorm:"default:generate_ulid()"`
	CheckID       *uuid.UUID `json:"check_id,omitempty"`
	CheckTime     *time.Time `json:"check_time,omitempty" time_format:"postgres_timestamp"`
	PlaybookRunID *uuid.UUID `json:"playbook_run_id,omitempty"`
	ConnectionID  uuid.UUID  `json:"connection_id,omitempty"`
	Path          string     `json:"path"`
	Filename      string     `json:"filename"`
	Size          int64      `json:"size"` // Size in bytes
	ContentType   string     `json:"content_type,omitempty"`
	Checksum      string     `json:"checksum"`
	CreatedAt     time.Time  `json:"created_at" yaml:"created_at" time_format:"postgres_timestamp"`
	UpdatedAt     time.Time  `json:"updated_at" yaml:"updated_at" time_format:"postgres_timestamp"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" yaml:"deleted_at,omitempty" time_format:"postgres_timestamp"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" yaml:"expires_at,omitempty" time_format:"postgres_timestamp"`
}
