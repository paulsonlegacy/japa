package entity

import "time"


// Plan defines the different subscription plans you offer
type Plan struct {
	ID               uint          `gorm:"type:tinyint;primaryKey"`

	Name             string       `gorm:"column:name;type:varchar(12);unique;not null"`   // e.g., "Basic", "Pro"
	Description      string       `gorm:"column:description;type:text"`   // A description of what this plan includes
	Price            float64      `gorm:"column:price;not null"`         // In your base currency (e.g., NGN)
	BillingCycle     string       `gorm:"column:billing_cycle;type:varchar(12);not null"` // "monthly", "yearly", "quarterly" etc.

	CreatedAt        time.Time
	UpdatedAt        time.Time

	PlanFeatures     []PlanFeature `gorm:"foreignKey:PlanID;references:ID"`
}
