// DB connection setup
package db

import (
	//"time"

	"japa/internal/config"
	"japa/internal/domain/entity"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


// Initialize GORM DB
func NewGormDB(cfg config.DBConfig) *gorm.DB {
	// DSN
	dsn := cfg.DBURL

	zap.S().Debugw("Opening database connection", "dburl", dsn)

	// Opening database pool
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("Failed to open database connection", zap.Error(err))
		panic("Database connection failed: " + err.Error())
	}

	// Get the generic database object sql.DB to configure connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		zap.L().Error("Failed to get sql.DB from gorm DB", zap.Error(err))
		panic("Failed to access DB internals: " + err.Error())
	}

	// Connection Pool Settings
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)
	sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTime)

	zap.L().Debug("Database connection established. Starting migration...")

	// Auto-migrate all models
	if err := gormDB.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.Subscription{},
		&entity.Purchase{},
		&entity.Plan{},
		&entity.PlanFeature{},
		&entity.Post{},
		&entity.Comment{},
		&entity.ScrapedPost{},
		//&entity.Reply{},
		//&entity.VisaFormInput{},
		&entity.VisaApplication{},
		&entity.Document{},
	); err != nil {
		zap.L().Error("Database migration failed", zap.Error(err))
		panic("Database migration failed: " + err.Error())
	}

	zap.L().Debug("Database migration completed successfully!")

	// Return *gorm.DB
	return gormDB
}
