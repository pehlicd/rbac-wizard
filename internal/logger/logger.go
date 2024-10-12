package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string, format string) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(levelFromString(level)).With().Timestamp().Logger()
	if format == "json" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else if format == "text" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})
	}

	return &logger
}

func levelFromString(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "off":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}
