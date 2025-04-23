package logging

import (
	"os"
	"japa/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


// Logger is your globally accessible logger instance
var Logger *zap.Logger


// InitLogger sets up and returns a configured zap logger based on your LoggingConfig.
// It also registers this logger as the global zap logger, so zap.L() will use it.
// Supports both custom logger (Logger.Info) and zap global (zap.L().Info).
func InitLogger(cfg config.LoggingConfig) *zap.Logger {
	var err error

	// Open general log file (INFO, DEBUG, etc.)
	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("could not open log file: " + err.Error())
	}

	// Open separate error log file (ERROR and above)
	errorLogFile, err := os.OpenFile(cfg.ErrorLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("could not open error log file: " + err.Error())
	}

	// Set logging level depending on the environment
	logLevel := zapcore.InfoLevel
	if cfg.EnvType == "LOCAL-DEV" {
		logLevel = zapcore.DebugLevel // more verbose during development
	}

	// Create a Tee core that writes:
	// - all logs >= logLevel to the main log file
	// - all logs >= ErrorLevel to the error log file
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(logFile),
			logLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(errorLogFile),
			zapcore.ErrorLevel,
		),
	)

	// Create the logger and attach caller info (file, line)
	logger := zap.New(core, zap.AddCaller())

	// Set the custom global logger so you can do logging.Logger.Info(...) everywhere
	Logger = logger

	// ALSO replace zap's global logger so zap.L().Info(...) works with this config
	zap.ReplaceGlobals(logger)

	return logger
}