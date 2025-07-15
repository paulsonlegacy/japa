package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)

// Visa application model
type VisaApplication struct {
	ID              string       `gorm:"type:char(26);primaryKey"` // Primary Key for the visa application
	UserID          string       `gorm:"column:user_id;not null"` // Who applied (ForeignKey to User)
	User            User         `gorm:"foreignKey:UserID"` // The user who is applying for the visa

	// Optional: provide form fields instead of uploading a form
	// JSON column for dynamic visa form input
	VisaFormInput   []byte        `gorm:"column:visa_form_input;type:json"` // Or use gorm.io/datatypes
	//VisaFormInput    *VisaFormInput `gorm:"foreignKey:VisaApplicationID;references:ID"` // One-to-One relationship (nullable)

	// Optional: signed visa uploaded form
	VisaFormURL     *string       `gorm:"column:visa_form_url;null"` // URL to the uploaded visa form (nullable)

	// Supporting documents
	Documents       []Document    `gorm:"foreignKey:VisaApplicationID;references:ID"` // Foreign Key for the documents

	// Assigned agent
	AgentID         *string       `gorm:"column:agent_id;null"` // Foreign Key to the assigned agent (nullable)
	Agent           *User         `gorm:"foreignKey:AgentID"` // The agent assigned to the application

	// Track status of the application
	Status          *string       `gorm:"column:status;null;default:'pending'"` // The status of the application (e.g., "pending", "submitted", "under review", "approved"  etc.)
	Feedback        *string       `gorm:"column:feedback;null"` // Feedback from the embassy or agent

	// Timestamps
	CreatedAt       time.Time    // Automatically set by GORM
	UpdatedAt       time.Time    // Automatically set by GORM
}

/*
{
  "user_id": "qwerttiitopppodakdkl",
  "visa_form_input": {     // can be null
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
  },
  "visa_form_url": "https://form.url" // can be null
}
*/