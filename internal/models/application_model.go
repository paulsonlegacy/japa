package models

import (
	"time"
)

type Application struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Country   string    `gorm:"size:100;not null"`
	VisaType  string    `gorm:"size:100;not null"`
	Status    string    `gorm:"default:pending"` // pending, approved, rejected
	Notes     string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
