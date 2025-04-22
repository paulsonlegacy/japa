// Loads env/config from .env or flags
package config

import "time"


// This gives you a singleton-style global access point to your config 
// across the app, without having to pass cfg *Config around manually
// i.e jwtSecret := config.Settings.JWT.JWTSecretKey
var Settings Config

// Config Types
type DatabaseConfig struct {
	DBURL        string
	DBDRIVER     string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
	MaxLifeTime  time.Duration
}


type JWTConfig struct {
	JWTSecretKey string
	Issuer       string
	Expiry       time.Duration
}

/*
1. EnvType: This tells your logger how “verbose” it should be.
You can use it to switch between environments like:

EnvType | Meaning | Typical Zap Level
LOCAL-DEV | Developer’s local machine | DebugLevel
STAGING | Pre-production test env | InfoLevel
PRODUCTION | Live user-facing system | Info or Error

Valid EnvType values
👉 LOCAL-DEV | STAGING | PRODUCTION | TEST

2. LogFilePath and ErrorLogFilePath
These are the paths where Zap will store logs (as files).
You're creating two log “streams”:

General logs → LogFilePath (e.g., japa.log)

Only error-level logs → ErrorLogFilePath (e.g., japa-errors.log)
*/
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


// Initialize configurations
func InitConfig(envFilePath string) *Config {
	loadEnvFile(envFilePath)
	cfg := &Config{
		Database: DatabaseConfig{
			DBURL:     getEnv("DBURL", ""),
			DBDRIVER:   getEnv("DBDRIVER", "mysql"),
			MaxOpenConns: 50,
			MaxIdleConns: 10,
			MaxIdleTime:  mustParseDuration("15m"),
			MaxLifeTime:  mustParseDuration("1h"),
		},
		JWT: JWTConfig{
			JWTSecretKey: getEnv("JWT_SECRET_KEY", ""),
			Issuer:       getEnv("JWT_ISSUER", "japa"),
			Expiry:       time.Hour * 24,
		},
		LoggingConfig: LoggingConfig{
			EnvType:          getEnv("ENV_TYPE", "LOCAL-DEV"),
			LogFilePath:      getEnv("LOG_FILE_PATH", "./logs/japa.log"),
			ErrorLogFilePath: getEnv("ERROR_LOG_FILE_PATH", "./logs/japa-errors.log"),
		},
	}
	Settings = *cfg
	return cfg
}

// Parse string to time.Duration
func mustParseDuration(val string) time.Duration {
	duration, err := time.ParseDuration(val)
	if err != nil {
		panic("invalid duration: " + val)
	}
	return duration
}