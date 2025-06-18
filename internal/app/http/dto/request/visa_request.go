package request

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type PersonalInfoRequest struct {
	PassportNumber   string    `json:"passport_number" validate:"required,len=9"` // common length for passport numbers
	PassportExpiry   string    `json:"passport_expiry" validate:"required"`
	ResidentialAddr  string    `json:"residential_address" validate:"required"`
	Nationality      string    `json:"nationality" validate:"required"`
	MaritalStatus    string    `json:"marital_status" validate:"required"`
	DateOfBirth      string    `json:"date_of_birth" validate:"required"`
}

type EmergencyContactRequest struct {
	EmergencyName     *string   `json:"emergency_name"`
	EmergencyPhone    *string   `json:"emergency_phone" validate:"omitempty,e164"` // international phone format (+...)
	EmergencyRelation *string   `json:"emergency_relation"`
}

type VisaFormInputRequest struct {
	Destination      string    `json:"destination" validate:"required,min=2"`
	VisaType         string    `json:"visa_type" validate:"required"`
	TravelDate       string    `json:"travel_date" validate:"required"`
	DurationOfStay   string    `json:"duration_of_stay" validate:"required"`
	Purpose          string    `json:"purpose" validate:"required"`
	HasBeenDenied    bool      `json:"has_been_denied"`

	// Personal details
	PersonalInfo      PersonalInfoRequest     `json:"personal_info" validate:"required"`

	// Emergency contact
	EmergencyContact  *EmergencyContactRequest `json:"emergency_contact"  validate:"omitempty,dive"`
}

type CreateVisaApplicationRequest struct {
	UserID  string    `json:"user_id" validate:"required,ulid"`

	// Optional fields
	VisaFormInput *VisaFormInputRequest `json:"visa_form_input"`
	VisaFormURL  *string `json:"visa_form_url" validate:"omitempty,url"`
}


func (req *CreateVisaApplicationRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	// Parse request body into req
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// Validate top-level fields
	if err := v.Struct(req); err != nil {
		return err
	}

	// If both VisaFormInput and VisaFormURL are nil
	if req.VisaFormInput == nil && req.VisaFormURL == nil {
		return errors.New("Both application input and form URL cannot be empty")
	}

	// If VisaFormInput is not nil, validate it separately
	if req.VisaFormInput != nil {
		if err := v.Struct(req.VisaFormInput); err != nil {
			return err
		}

		// Validate nested EmergencyContact if present
		if req.VisaFormInput.EmergencyContact != nil {
			if err := v.Struct(req.VisaFormInput.EmergencyContact); err != nil {
				return err
			}
		}
	}

	return nil
}


/*
{
  "user_id": "user_id",
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