package usecase

import (
	"context"
	"errors"
	//"fmt"
	"time"

	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"
	"japa/internal/domain/repository"
	"japa/internal/infrastructure/mail"
	"japa/internal/pkg"
	"japa/internal/util"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TYPES

// UserUsecase handles user-related business logic
type UserUsecase struct {
	Repo *repository.UserRepository
	DB   *gorm.DB
	Mail *mailer.SMTPMailer
}

// Mailer interface allows you to plug in real or mock email senders
type Mailer interface {
	SendWelcomeEmail(to string) error
}

// Initialize UserUsecase
func NewUserUsecase(repo *repository.UserRepository, db *gorm.DB, mailer *mailer.SMTPMailer) *UserUsecase {
	return &UserUsecase{Repo: repo, DB: db, Mail: mailer}
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

		// 2. Send welcome email in goroutine with timeout
		zap.L().Info("Forwarding welcome email..")
		emailData := mailer.WelcomeMail(user.Username)

		// Goroutine channel
		done := make(chan error, 1)

		go func() {
			// Simulate email sending
			err := usecase.Mail.SendViaSMTP(user.Email, emailData)
			done <- err
		}()

		select {
			case err := <-done:
				if err != nil {
					return err // rollback
				}
			case <-time.After(45 * time.Second):
				return errors.New("email send timeout") // rollback
		}

		// Everything succeeded
		return nil // commit
	})
}

// Logs in user based on credentials
func (us *UserUsecase) Login(account string, password string) (string, error) {
	// Find user by email or ussername 
	user, err := us.Repo.FindByEmailOrUsername(account)
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
