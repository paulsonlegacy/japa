// DB interaction logic using GORM
package repository

import (
	"context"
	"japa/internal/domain/entity"

	//"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// TYPES

// Repository method signatures
type UserRepositoryInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

// UserRepository to interface with DB
type UserRepository struct {
	DB *gorm.DB
}


// METHODS

// Initialize UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}


// Create user
func (ur *UserRepository) CreateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	return tx.WithContext(ctx).Create(user).Error
}


// Find user by email or username
func (ur *UserRepository) FindUserByEmailOrUsername(ctx context.Context, identifier string) (*entity.User, error) {
	var user entity.User
	err := ur.DB.
		WithContext(ctx).
		Where("email = ? OR username = ?", identifier, identifier).
		First(&user).Error
	return &user, err
}


// Find user by refresh token
func (ur *UserRepository) FindUserByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	var rt entity.RefreshToken
	if err := ur.DB.
		WithContext(ctx).
		Where("token = ?", refreshToken).
		First(&rt).Error; err != nil {
		return nil, err
	}

	var user entity.User
	if err := ur.DB.
		WithContext(ctx).
		Where("id = ?", rt.UserID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}


// Find user by id
func (ur *UserRepository) FindUserByID(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User
	if err := ur.DB.
		WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}


// Store refresh token
func (ur *UserRepository) SaveRefreshToken(ctx context.Context, refreshToken *entity.RefreshToken) error {
	return ur.DB.WithContext(ctx).Create(refreshToken).Error
}


// Delete refresh token
func (ur *UserRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	return ur.DB.WithContext(ctx).Where("token = ?", refreshToken).Delete(&entity.RefreshToken{}).Error
}