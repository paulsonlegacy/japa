// GORM user model struct
package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	FullName  string    `gorm:"size:100;not null"`
	Username  string    `gorm:"size:60;not null"`
	Email     string    `gorm:"size:100;unique;not null"`
	Password  string    `gorm:"not null"` // store hashed password
	Phone     string    `gorm:"size:20"`
	Role      string    `gorm:"default:user"` // can be "user" or "admin"
	CreatedAt time.Time
	UpdatedAt time.Time
	Applications []Application `gorm:"foreignKey:UserID"`
}