package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)


// Subscription tracks one subscription instance
type Subscription struct {
	ID          string     `gorm:"type:varchar(60);primaryKey"`

	UserID      string     `gorm:"column:user_id;type:varchar(60);not null"` // FK to User
	User        User       `gorm:"foreignKey:UserID"`

	PlanID      uint       `gorm:"column:plan_id;type:tinyint;not null"` // FK to Plan
	Plan        Plan       `gorm:"foreignKey:PlanID"`

	Status      string     `gorm:"column:status;type:varchar(60);not null;default:'active'"` // "active", "canceled", "expired"

	StartedAt   time.Time  `gorm:"column:started_at;not null"` // When subscription began
	ExpiresAt   time.Time  `gorm:"column:expires_at;not null"` // When it ends
	CanceledAt  *time.Time `gorm:"column:canceled_at"`         // Null unless user canceled early

	CreatedAt   time.Time
	UpdatedAt   time.Time
}