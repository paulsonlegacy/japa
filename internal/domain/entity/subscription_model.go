package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)


// Subscription tracks one subscription instance
type Subscription struct {
	ID          string     `gorm:"type:char(26);primaryKey"`

	UserID      string     `gorm:"column:user_id;not null"` // FK to User
	User        User       `gorm:"foreignKey:UserID"`

	PlanID      string     `gorm:"column:plan_id;not null"` // FK to Plan
	Plan        Plan       `gorm:"foreignKey:PlanID"`

	Status      string     `gorm:"column:status;not null;default:'active'"` // "active", "canceled", "expired"

	StartedAt   time.Time  `gorm:"column:started_at;not null"` // When subscription began
	ExpiresAt   time.Time  `gorm:"column:expires_at;not null"` // When it ends
	CanceledAt  *time.Time `gorm:"column:canceled_at"`         // Null unless user canceled early

	CreatedAt   time.Time
	UpdatedAt   time.Time
}