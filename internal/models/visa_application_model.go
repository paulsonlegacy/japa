package models

import (
	"time"
	"math/rand"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// VisaApplication represents an application request submitted by a user
// This holds core details needed by the agents to start visa processing

type VisaApplication struct {
	ID              ulid.ULID    `gorm:"type:char(26);primaryKey" json:"id"`
	UserID          ulid.ULID    `gorm:"not null" json:"user_id"` // Who applied
	User            User         `gorm:"foreignKey:UserID" json:"user"`

	Destination     string       `gorm:"not null" json:"destination"` // Country the user is applying to
	VisaType        string       `gorm:"not null" json:"visa_type"`   // e.g. Student, Work, Tourist
	TravelDate      time.Time    `gorm:"not null" json:"travel_date"`
	DurationOfStay  string       `gorm:"not null" json:"duration_of_stay"` // e.g. 3 months, 1 year
	Purpose         string       `gorm:"not null" json:"purpose"`
	HasBeenDenied   bool         `gorm:"not null;default:false" json:"has_been_denied"`  // Has user been denied visa before

	// Personal details
	PassportNumber  string       `gorm:"not null" json:"passport_number"`
	PassportExpiry  time.Time    `gorm:"not null" json:"passport_expiry"`
	ResidentialAddr string       `gorm:"not null" json:"residential_address"`
	Nationality     string       `gorm:"not null" json:"nationality"`
	MaritalStatus   string       `gorm:"not null" json:"marital_status"`
	DateOfBirth     time.Time    `gorm:"not null" json:"date_of_birth"`

	// Emergency contact details
	EmergencyName     string     `gorm:"not null" json:"emergency_name"`
	EmergencyPhone    string     `gorm:"not null" json:"emergency_phone"`
	EmergencyRelation string     `gorm:"not null" json:"emergency_relation"`

	// Track status of application
	Status          string       `gorm:"default:'Pending'" json:"status"` // e.g., "pending", "submitted", "in progress", "under review", "approved" etc.
	Notes           string       `gorm:"not null" json:"notes"`      // Internal note
	Feedback        string       `gorm:"not null" json:"feedback"`   // Embassy or agent feedback
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`

	// Supporting documents
	Documents       []Document   `gorm:"foreignKey:VisaApplicationID" json:"documents"`

	// Assigned agent
	AgentID         ulid.ULID    `gorm:"null" json:"agent_id"` 
	Agent           User         `gorm:"foreignKey:AgentID" json:"agent"`

	// Signed Visa Uploaded Form
	VisaFormURL     string       `gorm:"null" json:"visa_form_url"` // URL to visa form
}



// BeforeCreate hook runs before a new record is inserted into the DB.
// We use this to generate a ULID for the primary key.
func (v *VisaApplication) BeforeCreate(tx *gorm.DB) (err error) {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	v.ID = ulid.MustNew(ulid.Timestamp(t), entropy)
	return nil
}



// Visa application input placeholder
// For validation, data clarity, control & mapping
type VisaApplicationInput struct {
	UserID           string    `json:"user_id" binding:"required,ulid"`
	Destination      string    `json:"destination" binding:"required,min=2"`
	VisaType         string    `json:"visa_type" binding:"required"`
	TravelDate       string `json:"travel_date" binding:"required"`
	DurationOfStay   string    `json:"duration_of_stay" binding:"required"`
	Purpose          string    `json:"purpose" binding:"required"`
	HasBeenDenied    bool      `json:"has_been_denied"`

	// Personal details
	PassportNumber   string    `json:"passport_number" binding:"required,len=9"` // common length for passport numbers
	PassportExpiry   string `json:"passport_expiry" binding:"required"`
	ResidentialAddr  string    `json:"residential_address" binding:"required"`
	Nationality      string    `json:"nationality" binding:"required"`
	MaritalStatus    string    `json:"marital_status" binding:"required"`
	DateOfBirth      string `json:"date_of_birth" binding:"required"`

	// Emergency contact
	EmergencyName     string `json:"emergency_name" binding:"required"`
	EmergencyPhone    string `json:"emergency_phone" binding:"required,e164"` // international phone format (+...)
	EmergencyRelation string `json:"emergency_relation" binding:"required"`

	// Optional fields
	AgentID       string `json:"agent_id" binding:"omitempty,ulid"`
	VisaFormURL   string `json:"visa_form_url" binding:"omitempty,url"`
}


// Converts VisaApplicationInput to VisaApplication model
func ToVisaApplication(input *VisaApplicationInput) (VisaApplication, error) {
	// Type conversion
	userID, err := ulid.Parse(input.UserID)
	if err != nil {
		return VisaApplication{}, err
	}

	var agentID ulid.ULID
	if input.AgentID != "" {
		agentID, err = ulid.Parse(input.AgentID)
		if err != nil {
			return VisaApplication{}, err
		}
	}

	var travelDate time.Time
	if input.TravelDate != "" {
		travelDate, err = time.Parse("2006-01-02", input.TravelDate)
		if err != nil {
			return VisaApplication{}, err
		}
	}
	
	var passportExpiry time.Time
	if input.PassportExpiry != "" {
		passportExpiry, err = time.Parse("2006-01-02", input.PassportExpiry)
		if err != nil {
			return VisaApplication{}, err
		}
	}

	var dob time.Time
	if input.DateOfBirth != "" {
		dob, err = time.Parse("2006-01-02", input.DateOfBirth)
		if err != nil {
			return VisaApplication{}, err
		}
	}

	return VisaApplication{
		ID:              ulid.Make(), // or generate elsewhere
		UserID:          userID,
		Destination:     input.Destination,
		VisaType:        input.VisaType,
		TravelDate:      travelDate,
		DurationOfStay:  input.DurationOfStay,
		Purpose:         input.Purpose,
		HasBeenDenied:   input.HasBeenDenied,
		PassportNumber:  input.PassportNumber,
		PassportExpiry:  passportExpiry,
		ResidentialAddr: input.ResidentialAddr,
		Nationality:     input.Nationality,
		MaritalStatus:   input.MaritalStatus,
		DateOfBirth:     dob,
		EmergencyName:   input.EmergencyName,
		EmergencyPhone:  input.EmergencyPhone,
		EmergencyRelation: input.EmergencyRelation,
		Status:            "Pending",
		AgentID:           agentID,
		VisaFormURL:       input.VisaFormURL,
	}, nil
}
