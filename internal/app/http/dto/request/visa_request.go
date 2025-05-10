package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type CreateVisaApplicationRequest struct {
	UserID           string    `json:"user_id" validate:"required,ulid"`
	Destination      string    `json:"destination" validate:"required,min=2"`
	VisaType         string    `json:"visa_type" validate:"required"`
	TravelDate       string `json:"travel_date" validate:"required"`
	DurationOfStay   string    `json:"duration_of_stay" validate:"required"`
	Purpose          string    `json:"purpose" validate:"required"`
	HasBeenDenied    bool      `json:"has_been_denied"`

	// Personal details
	PassportNumber   string    `json:"passport_number" validate:"required,len=9"` // common length for passport numbers
	PassportExpiry   string `json:"passport_expiry" validate:"required"`
	ResidentialAddr  string    `json:"residential_address" validate:"required"`
	Nationality      string    `json:"nationality" validate:"required"`
	MaritalStatus    string    `json:"marital_status" validate:"required"`
	DateOfBirth      string `json:"date_of_birth" validate:"required"`

	// Emergency contact
	EmergencyName     *string `json:"emergency_name" validate:"required"`
	EmergencyPhone    *string `json:"emergency_phone" validate:"required,e164"` // international phone format (+...)
	EmergencyRelation *string `json:"emergency_relation" validate:"required"`

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


// Converts VisaApplicationreq to VisaApplication model
func ToVisaApplication(req *CreateVisaApplicationRequest) (*entity.VisaApplication, error) {
	// Type conversion
	userID, err := ulid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	var travelDate time.Time
	if req.TravelDate != "" {
		travelDate, err = time.Parse("2006-01-02", req.TravelDate)
		if err != nil {
			return nil, err
		}
	}
	
	var passportExpiry time.Time
	if req.PassportExpiry != "" {
		passportExpiry, err = time.Parse("2006-01-02", req.PassportExpiry)
		if err != nil {
			return nil, err
		}
	}

	var dob time.Time
	if req.DateOfBirth != "" {
		dob, err = time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, err
		}
	}

	return &entity.VisaApplication{
		ID:              util.NewULID(),
		UserID:          userID,
		Destination:     req.Destination,
		VisaType:        req.VisaType,
		TravelDate:      travelDate,
		DurationOfStay:  req.DurationOfStay,
		Purpose:         req.Purpose,
		HasBeenDenied:   req.HasBeenDenied,
		PassportNumber:  req.PassportNumber,
		PassportExpiry:  passportExpiry,
		ResidentialAddr: req.ResidentialAddr,
		Nationality:     req.Nationality,
		MaritalStatus:   req.MaritalStatus,
		DateOfBirth:     dob,
		EmergencyName:   req.EmergencyName,
		EmergencyPhone:  req.EmergencyPhone,
		EmergencyRelation: req.EmergencyRelation,
		VisaFormURL:       req.VisaFormURL,
	}, nil
}