package entity

import (
	"time"

	//"github.com/oklog/ulid/v2"
)


// Purchase records single payments for visa application service
type Purchase struct {
	ID                 string            `gorm:"type:char(26);primaryKey"`

	UserID             string            `gorm:"column:user_id;not null"`
	User               User              `gorm:"foreignKey:UserID"`

	VisaApplicationID  string            `gorm:"column:visa_application_id;not null"`
	VisaApplication    VisaApplication   `gorm:"foreignKey:VisaApplicationID"`

	PurchasedAt        time.Time         `gorm:"column:purchased_at;not null"` // When the purchase happened

	CreatedAt          time.Time
	UpdatedAt          time.Time
}
