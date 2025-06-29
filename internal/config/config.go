// Loads env/config from .env or flags
package config

import "time"

// This gives you a singleton-style global access point to your config
// across the app, without having to pass cfg *Config around manually
// i.e jwtSecret := config.Settings.JWT.JWTSecretKey
var Settings Config

// Config Types
type DBConfig struct {
	DBURL        string
	DBDRIVER     string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
	MaxLifeTime  time.Duration
}

type SMTPConfig struct {
	EMAIL_HOST     string
	EMAIL_PORT     int
	EMAIL_USERNAME string
	EMAIL_PASSWORD string
}

type EmailConfig struct {
	SMTPConfig SMTPConfig
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

type ServerConfig struct {
	ServerAddress           string
	ServerPort              string
	AuthorizationHeaderPath string
	TemplatesDir            string
	AssetsDir               string
}

type SiteConfig struct {
	SiteName   string
	SiteDomain string
	SiteEmail  string
	LogoURL    string
}

type Config struct {
	SiteConfig       SiteConfig
	ServerConfig     ServerConfig
	DBConfig         DBConfig
	EmailConfig      EmailConfig
	JWTConfig        JWTConfig
	LoggingConfig    LoggingConfig
}

// Initialize configurations
func InitConfig(envFilePath string) *Config {
	loadEnvFile(envFilePath)
	cfg := &Config{
		SiteConfig: SiteConfig{
			SiteName:   getEnv("SITE_NAME", ""),
			SiteDomain: getEnv("SITE_DOMAIN", ""),
			SiteEmail:  getEnv("SITE_EMAIL", "legacywebhub@gamil.com"),
			LogoURL:    getEnv("LOGO_URL", ""),
		},
		ServerConfig: ServerConfig{
			ServerAddress:           getEnv("SERVER_ADDRESS", ":8080"),
			AuthorizationHeaderPath: "Authorization",
			TemplatesDir:            getEnv("TEMPLATE_DIR", "templates"),
			AssetsDir:               getEnv("ASSET_DIR", "assets"),
		},
		DBConfig: DBConfig{
			DBURL:        getEnv("DBURL", ""),
			DBDRIVER:     getEnv("DBDRIVER", "mysql"),
			MaxOpenConns: 50,
			MaxIdleConns: 10,
			MaxIdleTime:  mustParseDuration("15m"),
			MaxLifeTime:  mustParseDuration("1h"),
		},
		EmailConfig: EmailConfig{
			SMTPConfig: SMTPConfig{
				EMAIL_HOST:     getEnv("EMAIL_HOST", ""),
				EMAIL_PORT:     getEnvInt("EMAIL_PORT", 465),
				EMAIL_USERNAME: getEnv("EMAIL_USERNAME", ""),
				EMAIL_PASSWORD: getEnv("EMAIL_PASSWORD", ""),
			},
		},
		JWTConfig: JWTConfig{
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
