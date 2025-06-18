package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)


type PersonalInfo struct {
	PassportNumber  string       `gorm:"column:passport_number;not null"` // User's passport number
	PassportExpiry  time.Time    `gorm:"column:passport_expiry;not null"` // Passport expiry date
	ResidentialAddr string       `gorm:"column:residential_address;not null"` // User's residential address
	Nationality     string       `gorm:"column:nationality;not null"` // User's nationality
	MaritalStatus   string       `gorm:"column:marital_status;not null"` // User's marital status
	DateOfBirth     time.Time    `gorm:"column:date_of_birth;not null"` // User's date of birth
}

type EmergencyContact struct {
	EmergencyName     *string     `gorm:"column:emergency_name;null"` // Emergency contact's name
	EmergencyPhone    *string     `gorm:"column:emergency_phone;null"` // Emergency contact's phone number
	EmergencyRelation *string     `gorm:"column:emergency_relationship;null"` // Relationship to the emergency contact
}

type VisaFormInput struct {
	ID              ulid.ULID    `gorm:"type:char(26);primaryKey"` // Primary Key for the visa application

	Destination     string       `gorm:"column:destination;not null"` // Country the user is applying to
	VisaType        string       `gorm:"column:visa_type;not null"` // Type of visa (e.g., Student, Work, Tourist)
	TravelDate      time.Time    `gorm:"column:travel_date;not null"` // Travel date for the visa application
	DurationOfStay  string       `gorm:"column:duration_of_stay;not null"` // Duration of stay (e.g., 3 months, 1 year)
	Purpose         string       `gorm:"column:purpose;not null"` // Purpose of the travel
	HasBeenDenied   bool         `gorm:"column:has_been_denied;not null;default:false"` // If the user has been denied a visa before
	
	// Personal details
	PersonalInfo       PersonalInfo      `gorm:"embedded"`

	// Emergency contact details
	EmergencyContact   *EmergencyContact  `gorm:"embedded"`

	// Others
	Note           *string       `gorm:"column:note;null"` // Internal note for the application
}