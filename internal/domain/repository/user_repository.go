// DB interaction logic using GORM
package repository

import (
	"japa/internal/domain/entity"
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
func (ur *UserRepository) Create(tx *gorm.DB, user *entity.User) error {
	return tx.Create(user).Error
}


// Find user by email or username
func (ur *UserRepository) FindByEmailOrUsername(identifier string) (*entity.User, error) {
	var user entity.User
	err := ur.DB.
		Where("email = ? OR username = ?", identifier, identifier).
		First(&user).Error
	return &user, err
}