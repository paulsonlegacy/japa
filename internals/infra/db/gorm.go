// DB connection setup
package db

import (
	"time"

	"japa/internals/config"
	//"japa/internals/user"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	GormDB *gorm.DB
}

func NewGormDB(cfg config.DatabaseConfig) *DB {
	dsn := cfg.DBURL

	zap.S().Debugw("Opening database connection", "dburl", dsn)

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
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	zap.L().Debug("Database connection established. Starting migration...")

	// Auto-migrate all models
	// if err := gormDB.AutoMigrate(&user.User{}); err != nil {
	// 	zap.L().Error("Database migration failed", zap.Error(err))
	// 	panic("Database migration failed: " + err.Error())
	// }

	//zap.L().Debug("Database migration completed successfully!")

	return &DB{GormDB: gormDB}
}
