package entity


type PlanFeature struct {
	ID           uint      `gorm:"type:tinyint;primaryKey"`

	PlanID       uint      `gorm:"column:plan_id;type:tinyint;not null"` // FK to Plan
	Plan         Plan      `gorm:"foreignKey:PlanID"`
	
	FeatureLabel string    `gorm:"column:feature_label;not null"` // i.e Blog access, premium travel insights
	FeatureValue string    `gorm:"column:feature_value;not null"` // Yes, Yes
}
