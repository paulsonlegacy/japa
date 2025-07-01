package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

//  refresh_tokens table
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    UserID    ulid.ULID `gorm:"not null;index"`
    Token     string    `gorm:"not null;uniqueIndex"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}