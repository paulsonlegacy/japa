package usecase

import (
	"time"
	"context"
	//"errors"

	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"
	"japa/internal/domain/repository"
	"japa/internal/pkg"
	"japa/internal/util"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/oklog/ulid/v2"
)


// UserUsecase handles user-related business logic
type VisaApplicationUsecase struct {
	Repo *repository.VisaRepository
	DB   *gorm.DB
	//Mailer Mailer
}

// METHODS

// Initialize UserUsecase
func NewVisaApplicationUsecase(repo *repository.VisaRepository, db *gorm.DB) *VisaApplicationUsecase {
	return &VisaApplicationUsecase{Repo: repo, DB: db}
}

// Creates a new visa application and sends a confirmation email
func (usecase *VisaApplicationUsecase) CreateApplication(ctx context.Context, req request.CreateVisaApplicationRequest) error {
	return usecase.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Type conversions
		userID, err := ulid.Parse(req.UserID)
		if err != nil {
			return err
		}

		var travelDate time.Time
		if req.TravelDate != "" {
			travelDate, err = time.Parse("2006-01-02", req.TravelDate)
			if err != nil {
				return err
			}
		}
		
		var passportExpiry time.Time
		if req.PersonalInfo.PassportExpiry != "" {
			passportExpiry, err = time.Parse("2006-01-02", req.PersonalInfo.PassportExpiry)
			if err != nil {
				return err
			}
		}

		var dob time.Time
		if req.PersonalInfo.DateOfBirth != "" {
			dob, err = time.Parse("2006-01-02", req.PersonalInfo.DateOfBirth)
			if err != nil {
				return err
			}
		}

		// 2. Save 
		zap.L().Info("Saving application to DB..")

		PersonalInfo := &entity.PersonalInfo{
			PassportNumber: req.PersonalInfo.PassportNumber,
			PassportExpiry: passportExpiry,
			ResidentialAddr: req.PersonalInfo.ResidentialAddr,
			Nationality: req.PersonalInfo.Nationality,
			MaritalStatus: req.PersonalInfo.MaritalStatus,
			DateOfBirth: dob,
		}

		EmergencyContact := &entity.EmergencyContact{
			EmergencyName: req.EmergencyContact.EmergencyName,
			EmergencyPhone: req.EmergencyContact.EmergencyPhone,
			EmergencyRelation: req.EmergencyContact.EmergencyRelation,
		}

		application := &entity.VisaApplication{
			ID:              util.NewULID(),
			UserID:          userID,
			Destination:     req.Destination,
			VisaType:        req.VisaType,
			TravelDate:      travelDate,
			DurationOfStay:  req.DurationOfStay,
			Purpose:         req.Purpose,
			HasBeenDenied:   req.HasBeenDenied,
			// Personal Info
			PersonalInfo: PersonalInfo,
			// Emergency Contact
			EmergencyContact: EmergencyContact,
			// Form URL
			VisaFormURL:       req.VisaFormURL,
		}


		if err := usecase.Repo.Create(tx, application); err != nil {
			return err // rollback
		}

		/*
			// 2. Send welcome email in goroutine with timeout
			zap.L().Info("Forwarding welcome email..")
			done := make(chan error, 1)

			go func() {
				// Simulate email sending
				err := mailer.SendWelcomeEmail(user)
				done <- err
			}()

			select {
			case err := <-done:
				if err != nil {
					return err // rollback
				}
			case <-time.After(5 * time.Second):
				return errors.New("email send timeout") // rollback
			}
		*/

		// Everything succeeded
		return nil // commit
	})
}