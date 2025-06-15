// DB interaction logic using GORM
package repository

import (
	"japa/internal/domain/entity"
	
	"gorm.io/gorm"
)

// TYPES

// Repository method signatures
type VisaRepositoryInterface interface {
	Create(visa *entity.VisaApplication) error
	FindByApplicationID(ID string) (*entity.VisaApplication, error)
}

// VisaRepository to interface with DB
type VisaRepository struct {
	DB *gorm.DB
}

