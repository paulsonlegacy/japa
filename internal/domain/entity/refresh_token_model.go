package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)

//  refresh_tokens table
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    UserID    string    `gorm:"type:varchar(60);not null;index"`
    Token     string    `gorm:"type:text;not null;uniqueIndex"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}