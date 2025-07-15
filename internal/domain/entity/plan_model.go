package entity

import "time"


// Plan defines the different subscription plans you offer
type Plan struct {
	ID               uint          `gorm:"type:int;primaryKey"`

	Name             string       `gorm:"column:name;not null"`          // e.g., "Basic", "Pro"
	Description      string       `gorm:"column:description;"`           // A description of what this plan includes
	Price            float64      `gorm:"column:price;not null"`         // In your base currency (e.g., NGN)
	BillingCycle     string       `gorm:"column:billing_cycle;not null"` // "monthly", "yearly", etc.

	CreatedAt        time.Time
	UpdatedAt        time.Time

	PlanFeatures     []PlanFeature `gorm:"foreignKey:PlanID;references:ID"`
}
