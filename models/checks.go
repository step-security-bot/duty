package models

import (
	"fmt"
	"time"

	"github.com/flanksource/duty/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CheckHealthStatus string

const (
	CheckStatusHealthy   = "healthy"
	CheckStatusUnhealthy = "unhealthy"
)

var CheckHealthStatuses = []CheckHealthStatus{
	CheckStatusHealthy,
	CheckStatusUnhealthy,
}

type Check struct {
	ID                 uuid.UUID           `json:"id" gorm:"default:generate_ulid()"`
	CanaryID           uuid.UUID           `json:"canary_id"`
	AgentID            uuid.UUID           `json:"agent_id,omitempty"`
	Spec               types.JSON          `json:"-"`
	Type               string              `json:"type"`
	Name               string              `json:"name"`
	Namespace          string              `json:"namespace"`
	Labels             types.JSONStringMap `json:"labels" gorm:"type:jsonstringmap"`
	Description        string              `json:"description,omitempty"`
	Status             CheckHealthStatus   `json:"status,omitempty"`
	Owner              string              `json:"owner,omitempty"`
	Severity           Severity            `json:"severity,omitempty"`
	Icon               string              `json:"icon,omitempty"`
	Transformed        bool                `json:"transformed,omitempty"`
	LastRuntime        *time.Time          `json:"last_runtime,omitempty"`
	NextRuntime        *time.Time          `json:"next_runtime,omitempty"`
	LastTransitionTime *time.Time          `json:"last_transition_time,omitempty"`
	CreatedAt          *time.Time          `json:"created_at,omitempty"`
	UpdatedAt          *time.Time          `json:"updated_at,omitempty" gorm:"autoUpdateTime:false"`
	DeletedAt          *time.Time          `json:"deleted_at,omitempty"`
	SilencedAt         *time.Time          `json:"silenced_at,omitempty"`

	// Auxiliary fields
	CanaryName   string        `json:"canary_name,omitempty" gorm:"-"`
	ComponentIDs []string      `json:"components,omitempty"  gorm:"-"` // Linked component ids
	Uptime       types.Uptime  `json:"uptime,omitempty"  gorm:"-"`
	Latency      types.Latency `json:"latency,omitempty"  gorm:"-"`
	Statuses     []CheckStatus `json:"checkStatuses,omitempty"  gorm:"-"`
	DisplayType  string        `json:"display_type,omitempty"  gorm:"-"`

	// These are calculated for the selected date range
	EarliestRuntime *time.Time `json:"earliestRuntime,omitempty" gorm:"-"`
	LatestRuntime   *time.Time `json:"latestRuntime,omitempty" gorm:"-"`
	TotalRuns       int        `json:"totalRuns,omitempty" gorm:"-"`
}

func (c Check) TableName() string {
	return "checks"
}

func (c Check) ToString() string {
	return fmt.Sprintf("%s-%s-%s", c.Name, c.Type, c.Description)
}

func (c Check) GetDescription() string {
	return c.Description
}

func (c Check) AsMap(removeFields ...string) map[string]any {
	return asMap(c, removeFields...)
}

type Checks []*Check

func (c Checks) Len() int {
	return len(c)
}

func (c Checks) Less(i, j int) bool {
	return c[i].ToString() < c[j].ToString()
}

func (c Checks) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Checks) Find(key string) *Check {
	for _, check := range c {
		if check.Name == key {
			return check
		}
	}
	return nil
}

type CheckStatus struct {
	CheckID   uuid.UUID `json:"check_id" gorm:"primaryKey"`
	Status    bool      `json:"status"`
	Invalid   bool      `json:"invalid,omitempty"`
	Time      string    `json:"time" gorm:"primaryKey"`
	Duration  int       `json:"duration"`
	Message   string    `json:"message,omitempty"`
	Error     string    `json:"error,omitempty"`
	Detail    any       `json:"-" gorm:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	// IsPushed when set to true indicates that the check status has been pushed to upstream.
	IsPushed bool `json:"is_pushed,omitempty"`
}

func (s CheckStatus) GetTime() (time.Time, error) {
	return time.Parse(time.DateTime, s.Time)
}

func (CheckStatus) TableName() string {
	return "check_statuses"
}

func (s CheckStatus) AsMap(removeFields ...string) map[string]any {
	return asMap(s, removeFields...)
}

// CheckStatusAggregate1h represents the `check_statuses_1h` table
type CheckStatusAggregate1h struct {
	CheckID   string    `gorm:"column:check_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Duration  int       `gorm:"column:duration"`
	Total     int       `gorm:"column:total"`
	Passed    int       `gorm:"column:passed"`
	Failed    int       `gorm:"column:failed"`
}

func (CheckStatusAggregate1h) TableName() string {
	return "check_statuses_1h"
}

// CheckStatusAggregate1d represents the `check_statuses_1d` table
type CheckStatusAggregate1d struct {
	CheckID   string    `gorm:"column:check_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Duration  int       `gorm:"column:duration"`
	Total     int       `gorm:"column:total"`
	Passed    int       `gorm:"column:passed"`
	Failed    int       `gorm:"column:failed"`
}

func (CheckStatusAggregate1d) TableName() string {
	return "check_statuses_1d"
}

// CheckSummary represents the `check_summary` view
type CheckSummary struct {
	ID                 uuid.UUID           `json:"id"`
	CanaryID           uuid.UUID           `json:"canary_id"`
	CanaryName         string              `json:"canary_name"`
	CanaryNamespace    string              `json:"canary_namespace"`
	Description        string              `json:"description,omitempty"`
	Icon               string              `json:"icon,omitempty"`
	Labels             types.JSONStringMap `json:"labels"`
	LastTransitionTime *time.Time          `json:"last_transition_time,omitempty"`
	Latency            types.Latency       `json:"latency,omitempty"`
	Name               string              `json:"name"`
	Namespace          string              `json:"namespace"`
	Owner              string              `json:"owner,omitempty"`
	Severity           string              `json:"severity,omitempty"`
	Status             string              `json:"status"`
	Type               string              `json:"type"`
	Uptime             types.Uptime        `json:"uptime,omitempty"`
	LastRuntime        *time.Time          `json:"last_runtime,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	DeletedAt          *time.Time          `json:"deleted_at,omitempty"`
	SilencedAt         *time.Time          `json:"silenced_at,omitempty"`
}

func (t *CheckSummary) TableName() string {
	return "check_summary"
}

func (t CheckSummary) AsMap(removeFields ...string) map[string]any {
	return asMap(t, removeFields...)
}

type CheckConfigRelationship struct {
	ConfigID   uuid.UUID  `json:"config_id,omitempty"`
	CheckID    uuid.UUID  `json:"check_id,omitempty"`
	CanaryID   uuid.UUID  `json:"canary_id,omitempty"`
	SelectorID string     `json:"selector_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

func (c *CheckConfigRelationship) Save(db *gorm.DB) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "canary_id"}, {Name: "check_id"}, {Name: "config_id"}, {Name: "selector_id"}},
		UpdateAll: true,
	}).Create(c).Error
}

func (CheckConfigRelationship) TableName() string {
	return "check_config_relationships"
}
