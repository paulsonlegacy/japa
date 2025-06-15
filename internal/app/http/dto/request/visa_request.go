package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type PersonalInfo struct {
	PassportNumber   string    `json:"passport_number" validate:"required,len=9"` // common length for passport numbers
	PassportExpiry   string    `json:"passport_expiry" validate:"required"`
	ResidentialAddr  string    `json:"residential_address" validate:"required"`
	Nationality      string    `json:"nationality" validate:"required"`
	MaritalStatus    string    `json:"marital_status" validate:"required"`
	DateOfBirth      string    `json:"date_of_birth" validate:"required"`
}

type EmergencyContact struct {
	EmergencyName     *string   `json:"emergency_name" validate:"required"`
	EmergencyPhone    *string   `json:"emergency_phone" validate:"required,e164"` // international phone format (+...)
	EmergencyRelation *string   `json:"emergency_relation" validate:"required"`
}

type CreateVisaApplicationRequest struct {
	UserID           string    `json:"user_id" validate:"required,ulid"`
	Destination      string    `json:"destination" validate:"required,min=2"`
	VisaType         string    `json:"visa_type" validate:"required"`
	TravelDate       string    `json:"travel_date" validate:"required"`
	DurationOfStay   string    `json:"duration_of_stay" validate:"required"`
	Purpose          string    `json:"purpose" validate:"required"`
	HasBeenDenied    bool      `json:"has_been_denied"`

	// Personal details
	PersonalInfo      PersonalInfo     `json:"personal_info" validate:"required"`

	// Emergency contact
	EmergencyContact  EmergencyContact `json:"emergency_contact"`

	// Optional fields
	VisaFormURL   *string `json:"visa_form_url" validate:"omitempty,url"`
}


func (req *CreateVisaApplicationRequest) Bind(c *fiber.Ctx, v validator.Validate) error {
	// Parse request body into req
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// Validate request struct
	if err := v.Struct(req); err != nil {
		return err
	}

	return nil
}