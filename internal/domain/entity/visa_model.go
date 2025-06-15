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

// Visa application model
type VisaApplication struct {
	ID              ulid.ULID    `gorm:"type:char(26);primaryKey"` // Primary Key for the visa application
	UserID          ulid.ULID    `gorm:"column:user_id;not null"` // Who applied (ForeignKey to User)
	User            User         `gorm:"foreignKey:UserID"` // The user who is applying for the visa

	Destination     string       `gorm:"column:destination;not null"` // Country the user is applying to
	VisaType        string       `gorm:"column:visa_type;not null"` // Type of visa (e.g., Student, Work, Tourist)
	TravelDate      time.Time    `gorm:"column:travel_date;not null"` // Travel date for the visa application
	DurationOfStay  string       `gorm:"column:duration_of_stay;not null"` // Duration of stay (e.g., 3 months, 1 year)
	Purpose         string       `gorm:"column:purpose;not null"` // Purpose of the travel
	HasBeenDenied   bool         `gorm:"column:has_been_denied;not null;default:false"` // If the user has been denied a visa before

	// Personal details
	PersonalInfo       *PersonalInfo      `gorm:"embedded"`

	// Emergency contact details
	EmergencyContact   *EmergencyContact  `gorm:"embedded"`

	// Others
	Note           *string       `gorm:"column:note;null"` // Internal note for the application

	// Supporting documents
	Documents       []Document   `gorm:"foreignKey:VisaApplicationID"` // Foreign Key for the documents

	// Signed Visa Uploaded Form
	VisaFormURL     *string       `gorm:"null"` // URL to the uploaded visa form (nullable)

	// Assigned agent
	AgentID         *ulid.ULID    `gorm:"column:agent_id;null"` // Foreign Key to the assigned agent (nullable)
	Agent           *User         `gorm:"foreignKey:AgentID"` // The agent assigned to the application

	// Track status of the application
	Status          string       `gorm:"column:status;not null;default:'Pending'"` // The status of the application (e.g., "Pending", "Submitted", "Under Review", "Approved"  etc.)
	Feedback        *string       `gorm:"column:feedback;null"` // Feedback from the embassy or agent

	// Timestamps
	CreatedAt       time.Time    // Automatically set by GORM
	UpdatedAt       time.Time    // Automatically set by GORM
}

/*
{
  "destination": "Germany",
  "visa_type": "Student",
  "travel_date": "2025-09-01T00:00:00Z",
  "duration_of_stay": "2 years",
  "purpose": "University Education",
  "has_been_denied": false,
  "personal_info": {
    "passport_number": "A12345678",
    "passport_expiry": "2029-12-31T00:00:00Z",
    "residential_address": "123 Lagos Street, Abuja",
    "nationality": "Nigerian",
    "marital_status": "Single",
    "date_of_birth": "2001-07-15T00:00:00Z"
  },
  "emergency_contact": {
    "emergency_name": "John Doe",
    "emergency_phone": "+2348012345678",
    "emergency_relationship": "Brother"
  }
}
*/