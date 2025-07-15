package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)

// Document stores user-uploaded documents per application
// Flexible structure to support varying requirements across visa types/countries
type Document struct {
	ID                string    `gorm:"type:varchar(60);primaryKey"`
	VisaApplicationID string    `gorm:"type:varchar(60);not null"`
	VisaApplication   VisaApplication `gorm:"foreignKey:VisaApplicationID"`

	FileType          string    `gorm:"type:varchar(60);not null"` // e.g. "passport_photo", "bank_statement", "signed_form"
	FilePath          string    `gorm:"type:text;not null"` // where it's stored locally or cloud URL
	UploadedAt        time.Time `gorm:"autoCreateTime"`
}
