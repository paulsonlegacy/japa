package models

import (
	"time"
	"math/rand"

	"golang.org/x/crypto/bcrypt"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// User represents the basic user profile in the system
// This could later be extended to include verification, profile image, etc.
type User struct {
	ID        ulid.ULID `gorm:"type:char(26);primaryKey"`
	FirstName  string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Phone     string    `gorm:"not null"`
	Password  string    `gorm:"not null"` // hashed password
	Role      string    `gorm:"not null"` // admin, user, agent
	CreatedAt time.Time
	UpdatedAt time.Time
	VisaApplications []VisaApplication `gorm:"foreignKey:UserID"`
}

// BeforeCreate hook runs before a new record is inserted into the DB.
// We use this to generate a ULID for the primary key.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	u.ID = ulid.MustNew(ulid.Timestamp(t), entropy)
	return nil
}



// User input placeholder
// For validation, data clarity, control & mapping
type CreateUserInput struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,e164"` // or use `min/max` if not using `e164`
	Password  string `json:"password" binding:"required,min=8"` // plain password; hash before saving
	Role      string `json:"role" binding:"required,oneof=admin user agent"` // customize roles as needed
}


// ToUser converts CreateUserInput to the User model
func ToUser(input *CreateUserInput) (User, error) {
	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	return User{
		ID:        ulid.Make(), // Generates a new ULID for the user
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}