package zerologger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"vk-internship/internal/config"
	"vk-internship/internal/logger"
)

type Logger struct {
	zerolog.Logger
}

func New(cfg *config.LoggerConfig) *Logger {
	logLevel := parseLogLevel(cfg.Level)
	zerolog.SetGlobalLevel(logLevel)

	var output io.Writer
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: zerolog.TimeFormatUnix,
		}
	} else {
		output = os.Stdout
	}

	base := zerolog.New(output).With().Timestamp().CallerWithSkipFrameCount(3).Logger()

	return &Logger{base}
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
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
	default:
		return zerolog.InfoLevel
	}
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}

func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.Logger.Warn().Msg(msg)
}

func (l *Logger) Error(err error, msg string) {
	l.Logger.Error().Err(err).Msg(msg)
}

func (l *Logger) Fatal(err error, msg string) {
	l.Logger.Fatal().Err(err).Msg(msg)
}

func (l *Logger) With(fields map[string]interface{}) logger.Logger {
	return &Logger{l.Logger.With().Fields(fields).Logger()}
}

func (l *Logger) Component(name string) logger.Logger {
	return &Logger{l.Logger.With().Str("component", name).Logger()}
}
