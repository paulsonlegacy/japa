// Business logic (e.g., login, profile update)
package services

import (
	//"time"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"japa/internal/models"
	"japa/internal/repository"
	"japa/internal/utils"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

// TYPES

// UserService handles user-related business logic
type UserService struct {
	Repo *repository.UserRepository
	DB   *gorm.DB
	//Mailer Mailer
}

// Mailer interface allows you to plug in real or mock email senders
type Mailer interface {
	SendWelcomeEmail(to string) error
}


// METHODS

// Initialize UserService
func NewUserService(repo *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{Repo: repo, DB: db}
}

// Registers a new user and sends a welcome email
func (us *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	return us.DB.Transaction(func (tx *gorm.DB) error {
		// 1. Save user
		zap.L().Info("Saving user to DB")

		user.Password = utils.HashPassword(user.Password)
		if err := us.Repo.Create(tx, user); err != nil {
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
func (us *UserService) Login(email, password string) (string, error) {
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
	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}