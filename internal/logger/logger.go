package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup configures the global logger
func Setup(level string, jsonOutput bool) {
	// Set the global log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Configure output format
	if jsonOutput {
		log.Logger = log.Output(os.Stdout)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	// Add caller information in debug mode
	if logLevel <= zerolog.DebugLevel {
		log.Logger = log.With().Caller().Logger()
	}
}

// WithContext returns a logger with context fields
func WithContext(fields map[string]interface{}) zerolog.Logger {
	logger := log.With()
	for k, v := range fields {
		logger = logger.Interface(k, v)
	}
	return logger.Logger()
}
