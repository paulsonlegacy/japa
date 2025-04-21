// Loads env/config from .env or flags
package config

import "time"

var Settings Config

type DatabaseConfig struct {
	DBURL        string
	DBDRIVER     string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}


type JWTConfig struct {
	JWTSecretKey string
	Issuer       string
	Expiry       time.Duration
}


type LoggingConfig struct {
	EnvType          string
	LogFilePath      string
	ErrorLogFilePath string
}

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	LoggingConfig LoggingConfig
}


func Init(envFilePath string) *Config {
	loadEnvFile(envFilePath)
	cfg := &Config{
		Database: DatabaseConfig{
			DBURL:     getEnv("BDURL", ""),
			DBDRIVER:   getEnv("DBDRIVER", "mysql"),
		},
		JWT: JWTConfig{
			JWTSecretKey: getEnv("JWT_SECRET_KEY", ""),
			Issuer:       getEnv("JWT_ISSUER", "japa"),
			Expiry:       time.Hour * 24,
		},
		LoggingConfig: LoggingConfig{
			EnvType:          getEnv("ENV_TYPE", "LOCAL-DEV"),
			LogFilePath:      getEnv("LOG_FILE_PATH", "./japa.log"),
			ErrorLogFilePath: getEnv("ERROR_LOG_FILE_PATH", "./japa-errors.log"),
		},
	}
	Settings = *cfg
	return cfg
}
