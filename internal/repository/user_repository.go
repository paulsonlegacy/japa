// DB interaction logic using GORM
package repository

import (
	"japa/internal/models"
	"gorm.io/gorm"
)

// TYPES

// Repository method signatures
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
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
func (ur *UserRepository) Create(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}


// Find user by email
func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}