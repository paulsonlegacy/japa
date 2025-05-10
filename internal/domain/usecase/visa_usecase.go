package usecase

import (
	"time"
	"context"
	"errors"

	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"
	"japa/internal/domain/repository"
	"japa/internal/pkg"
	"japa/internal/util"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)


// UserUsecase handles user-related business logic
type VisaApplicationUsecase struct {
	Repo *repository.UserRepository
	DB   *gorm.DB
	//Mailer Mailer
}

// METHODS

// Initialize UserUsecase
func NewVisaApplicationUsecase(repo *repository.UserRepository, db *gorm.DB) *UserUsecase {
	return &UserUsecase{Repo: repo, DB: db}
}

// Registers a new user and sends a welcome email
func (usecase *UserUsecase) RegisterUser(ctx context.Context, req request.CreateUserRequest) error {
	return usecase.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Save user
		zap.L().Info("Saving user to DB..")

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := &entity.User{
			ID:        util.NewULID(),
			FullName:  req.FullName,
			Username:  req.Username,
			Email:     req.Email,
			Phone:     req.Phone,
			Password:  string(hashedPassword),
			Role:      req.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}


		if err := usecase.Repo.Create(tx, user); err != nil {
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

// Logs in user based on credentials
func (us *UserUsecase) Login(email, password string) (string, error) {
	// Find user by email
	user, err := us.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Confirm password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT Token
	token, err := pkg.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
