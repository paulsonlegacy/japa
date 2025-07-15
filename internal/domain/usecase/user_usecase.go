package usecase

import (
	"context"
	"errors"
	"fmt"

	//"fmt"
	"time"

	"japa/internal/app/http/dto/request"
	"japa/internal/config"
	"japa/internal/domain/entity"
	"japa/internal/domain/repository"
	"japa/internal/infrastructure/mail"
	"japa/internal/pkg"

	//"japa/internal/util"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TYPES

// UserUsecase handles user-related business logic
type UserUsecase struct {
	JWTConfig  config.JWTConfig 
	Repo      *repository.UserRepository
	DB        *gorm.DB
	Mailer    *mailer.ResponsiveMailer
}

// Initialize UserUsecase
func NewUserUsecase(jwtConfig config.JWTConfig, repo *repository.UserRepository, db *gorm.DB, mailer *mailer.ResponsiveMailer) *UserUsecase {
	return &UserUsecase{JWTConfig: jwtConfig, Repo: repo, DB: db, Mailer: mailer}
}

// Registers a new user and sends a welcome email
func (usecase *UserUsecase) RegisterUser(ctx context.Context, req request.CreateUserRequest) error {
	return usecase.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Save user
		zap.L().Info("Saving user to DB..")

		// Hash password
		hashedPassword := pkg.HashAndEncodeArgon2(req.Password, 32)

		// Generating user id
		userID := ulid.Make().String()

		user := &entity.User{
			ID:        userID,
			FullName:  req.FullName,
			Username:  req.Username,
			Email:     req.Email,
			Phone:     req.Phone,
			Password:  string(hashedPassword),
			Role:      req.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}


		if err := usecase.Repo.CreateUser(ctx, tx, user); err != nil {
			return err // rollback
		}

		// 2. Send welcome email in goroutine with timeout
		//zap.L().Info("Forwarding welcome email..")
		emailData := mailer.WelcomeMail(user.Username)

		err := usecase.Mailer.Send(user.Email, emailData)
		if err != nil {
			fmt.Println("Welcome mail failed to send")
			return err // rollback
		}

		// Everything succeeded
		return nil // commit
	})
}

// Logs in user based on credentials
func (us *UserUsecase) LoginUser(ctx context.Context, account string, password string) (string, string, error) {
	/*
	// Note that if only one device is to be logged in at a time,
	// Delete previous user refresh tokens before saving a new one,
	// That way  only one device can be signed in at a time
	*/
	
	// Find user
	user, err := us.Repo.FindUserByEmailOrUsername(ctx, account)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Verify password
	if !pkg.Compare(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate access JWT
	accessToken, err := pkg.GenerateJWT(user, us.JWTConfig)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := pkg.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Saving refresh token to DB
	newToken := &entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	// Store refresh token in DB if you like:
	if err := us.Repo.SaveRefreshToken(ctx, newToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Refresh token provider
func (usecase *UserUsecase) Logout(ctx context.Context, refreshToken string) error {
	if err := usecase.Repo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

func (usecase *UserUsecase) GetRefreshToken(ctx context.Context, refreshToken string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	if err := usecase.DB.WithContext(ctx).First(&token, "token = ?", refreshToken).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &token, nil
}