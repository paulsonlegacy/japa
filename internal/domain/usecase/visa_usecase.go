package usecase

import (
	"time"
	"context"
	//"errors"
	"encoding/json"

	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"
	"japa/internal/domain/repository"
	//"japa/internal/util"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/oklog/ulid/v2"
)


// UserUsecase handles user-related business logic
type VisaUsecase struct {
	Repo *repository.VisaRepository
	DB   *gorm.DB
	//Mailer Mailer
}

// METHODS

// Initialize UserUsecase
func NewVisaUsecase(repo *repository.VisaRepository, db *gorm.DB) *VisaUsecase {
	return &VisaUsecase{Repo: repo, DB: db}
}

// Creates a new visa application and sends a confirmation email
func (usecase *VisaUsecase) CreateVisaApplication(ctx context.Context, req request.CreateVisaApplicationRequest) error {
	return usecase.DB.Transaction(func(tx *gorm.DB) error {
		// Application placeholder
		var application = &entity.VisaApplication{}

		// 1. Type conversions
		var userID string 
		parsedUserID, err := ulid.Parse(req.UserID)
		if err != nil {
			return err
		}
		userID = parsedUserID.String()

		// 2. Submitting by given input fields
		if req.VisaFormInput != nil {
			var travelDate time.Time
			if req.VisaFormInput.TravelDate != "" {
				travelDate, err = time.Parse("2006-01-02", req.VisaFormInput.TravelDate)
				if err != nil {
					return err
				}
			}

			var passportExpiry time.Time
			if req.VisaFormInput.PersonalInfo.PassportExpiry != "" {
				passportExpiry, err = time.Parse("2006-01-02", req.VisaFormInput.PersonalInfo.PassportExpiry)
				if err != nil {
					return err
				}
			}

			var dob time.Time
			if req.VisaFormInput.PersonalInfo.DateOfBirth != "" {
				dob, err = time.Parse("2006-01-02", req.VisaFormInput.PersonalInfo.DateOfBirth)
				if err != nil {
					return err
				}
			}

			visaFormInput := &entity.VisaFormInput{
				Destination:     req.VisaFormInput.Destination,
				VisaType:        req.VisaFormInput.VisaType,
				TravelDate:      travelDate,
				DurationOfStay:  req.VisaFormInput.DurationOfStay,
				Purpose:         req.VisaFormInput.Purpose,
				HasBeenDenied:   req.VisaFormInput.HasBeenDenied,
				PersonalInfo:    entity.PersonalInfo{
					PassportNumber:   req.VisaFormInput.PersonalInfo.PassportNumber,
					PassportExpiry:   passportExpiry,
					ResidentialAddr:  req.VisaFormInput.PersonalInfo.ResidentialAddr,
					Nationality:      req.VisaFormInput.PersonalInfo.Nationality,
					MaritalStatus:    req.VisaFormInput.PersonalInfo.MaritalStatus,
					DateOfBirth:      dob,
				},
				EmergencyContact: (*entity.EmergencyContact)(req.VisaFormInput.EmergencyContact),
			}

			jsonVisaFormInput, err := json.Marshal(visaFormInput)
			if err != nil {
				return err
			}

			application = &entity.VisaApplication{
				ID:                ulid.Make().String(),
				UserID:            userID,
				VisaFormInput:     jsonVisaFormInput,
				VisaFormURL:       req.VisaFormURL,
			}
		} else {
			jsonVisaFormInput, err := json.Marshal(req.VisaFormInput)
			if err != nil {
				return err
			}

			application = &entity.VisaApplication{
				ID:                ulid.Make().String(),
				UserID:            userID,
				VisaFormInput:     jsonVisaFormInput,
				VisaFormURL:       req.VisaFormURL,
			}
		}


		// 2. Save 
		zap.L().Info("Saving application to DB..")

		if err := usecase.Repo.Create(tx, application); err != nil {
			return err // rollback
		}

		// Everything succeeded
		return nil // commit
	})
}