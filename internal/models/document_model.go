package models

import (
	"time"
	"math/rand"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Document stores user-uploaded documents per application
// Flexible structure to support varying requirements across visa types/countries
type Document struct {
	ID                ulid.ULID `gorm:"type:char(26);primaryKey"`
	VisaApplicationID ulid.ULID `gorm:"not null"`
	VisaApplication   VisaApplication `gorm:"foreignKey:VisaApplicationID"`

	FileType      string    `gorm:"not null"` // e.g. "passport_photo", "bank_statement", "signed_form"
	FilePath      string    `gorm:"not null"` // where it's stored locally or cloud URL
	UploadedAt    time.Time `gorm:"autoCreateTime"`
}

// BeforeCreate hook runs before a new record is inserted into the DB.
// We use this to generate a ULID for the primary key.
func (d *Document) BeforeCreate(tx *gorm.DB) (err error) {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	d.ID = ulid.MustNew(ulid.Timestamp(t), entropy)
	return nil
}



// Document placeholder
// For validation, data clarity, control & mapping
type DocumentInput struct {
	VisaApplicationID string `json:"visa_application_id" binding:"required,len=26"` // ULID is 26 chars
	FileType          string `json:"file_type" binding:"required,oneof=passport_photo bank_statement signed_form"`
	FilePath          string `json:"file_path" binding:"required,url"`
}


// Converts DocumentInput to Document model
func ToDocument(input *DocumentInput) (Document, error) {
	// Parse VisaApplicationID
	visaAppID, err := ulid.Parse(input.VisaApplicationID)
	if err != nil {
		return Document{}, err
	}

	return Document{
		ID:                ulid.Make(), // generate a new ULID for the document
		VisaApplicationID: visaAppID,
		FileType:          input.FileType,
		FilePath:          input.FilePath,
	}, nil
}
