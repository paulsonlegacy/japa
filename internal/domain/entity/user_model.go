package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)

// User model
type User struct {
	ID                string    `gorm:"type:char(26);primaryKey"`
	FullName          string    `gorm:"column:full_name;type:varchar(100);not null"`                // explicit column name
	Username          string    `gorm:"column:username;type:varchar(60);uniqueIndex;not null"`     // unique
	Email             string    `gorm:"column:email;type:varchar(120);uniqueIndex;not null"`        // unique
	Phone             string    `gorm:"column:phone;type:varchar(30);uniqueIndex;not null"`        // unique
	Password          string    `gorm:"column:password;type:varchar(255);not null"`                 // hashed password
	Role              string    `gorm:"column:role;type:varchar(12);not null;default:user"`        // user, agent, admin, superadmin etc.
	BannedUntil       *time.Time `gorm:"column:banned_until;default:null"`
	BanReason         *string    `gorm:"column:ban_reason;type:varchar(200);default:null"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`         // GORM auto timestamps
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime"`         // GORM auto timestamps

	// Subscription relation
	Subscriptions     []Subscription `gorm:"foreignKey:UserID;references:ID"`
	// Visa application relation
	VisaApplications  []VisaApplication `gorm:"foreignKey:UserID;references:ID"`
}