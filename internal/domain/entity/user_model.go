package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// User model
type User struct {
	ID                ulid.ULID `gorm:"type:char(26);primaryKey"`
	FullName          string    `gorm:"column:full_name;not null"`                // explicit column name
	Username          string    `gorm:"column:username;uniqueIndex;not null"`     // unique
	Email             string    `gorm:"column:email;uniqueIndex;not null"`        // unique
	Phone             string    `gorm:"column:phone;uniqueIndex;not null"`        // unique
	Password          string    `gorm:"column:password;not null"`                 // hashed password
	Role              string    `gorm:"column:role;not null;default:user"`        // user, agent, etc.
	BannedUntil       *time.Time `gorm:"column:banned_until;default:null"`
	BanReason         *string    `gorm:"column:ban_reason;default:null"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`         // GORM auto timestamps
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime"`         // GORM auto timestamps

	// Subscription relation
	Subscriptions     []Subscription `gorm:"foreignKey:UserID;references:ID"`
	// Visa application relation
	VisaApplications  []VisaApplication `gorm:"foreignKey:UserID;references:ID"`
}