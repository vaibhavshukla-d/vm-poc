package logger

import (
	"vm/pkg/cinterface"
	"vm/pkg/constants"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	configmanager "vm/pkg/config-manager"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var zeroSinLogger *zerolog.Logger

type zeroLogger struct {
	cfg    *configmanager.ApplicationConfigModal
	logger *zerolog.Logger
}

func NewLogger(cfg *configmanager.Config) cinterface.Logger {
	return newZeroLogger(cfg)
}

var zeroLogLevelMapping = map[string]zerolog.Level{
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"fatal":    zerolog.FatalLevel,
	"disabled": zerolog.Disabled,
	"trace":    zerolog.TraceLevel,
}

func newZeroLogger(cfg *configmanager.Config) *zeroLogger {
	logger := &zeroLogger{cfg: cfg.App}
	logger.Init()
	return logger
}

func (l *zeroLogger) getLogLevel() zerolog.Level {
	level, exists := zeroLogLevelMapping[l.cfg.Application.Log.Level]
	if !exists {
		return zerolog.DebugLevel
	}
	return level
}

func (l *zeroLogger) Init() {

	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		// Create multi-writer for multiple outputs
		var writers []io.Writer

		// Add console writer if enabled
		if l.cfg.Application.Log.EnableConsole {
			consoleWriter := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
				NoColor:    false,
			}
			writers = append(writers, consoleWriter)
		}

		if l.cfg.Application.Log.Level != "disabled" && l.cfg.Application.Log.EnableFile {
			logDir := filepath.Join(l.cfg.Application.Log.FilePath, l.cfg.Application.Application.Name, "internal", time.Now().Format("2006-01-02"))

			// Ensure the directory exists
			if err := os.MkdirAll(logDir, 0755); err != nil {
				panic(fmt.Sprintf("failed to create log directory: %v", err))
			}

			// Create full log file path
			fileName := fmt.Sprintf("%s-%s.%s",
				l.cfg.Application.Log.FileName,
				uuid.New(),
				"log")
			fullPath := filepath.Join(logDir, fileName)

			// Open or create the log file
			file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				panic(fmt.Sprintf("could not open log file: %v", err))
			}

			writers = append(writers, file)
		}

		// Add file writer if enabled

		// If no writers are enabled, default to console
		if len(writers) == 0 {
			consoleWriter := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
				NoColor:    false,
			}
			writers = append(writers, consoleWriter)
		}

		// Create multi-writer
		mw := io.MultiWriter(writers...)

		// Create logger
		var logger = zerolog.New(mw).
			With().
			Timestamp().
			Str("AppName", l.cfg.Application.Application.Name).
			Str("LoggerName", "Zerolog").
			Logger()

		// Set log level
		zerolog.SetGlobalLevel(l.getLogLevel())
		zeroSinLogger = &logger
	})
	l.logger = zeroSinLogger
}

func (l *zeroLogger) Debug(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {

	l.logger.
		Debug().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Debugf(template string, args ...interface{}) {
	l.logger.
		Debug().
		Msgf(template, args...)
}

func (l *zeroLogger) Info(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {

	l.logger.
		Info().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Infof(template string, args ...interface{}) {
	l.logger.
		Info().
		Msgf(template, args...)
}

func (l *zeroLogger) Warn(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {

	l.logger.
		Warn().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Warnf(template string, args ...interface{}) {
	l.logger.
		Warn().
		Msgf(template, args...)
}

func (l *zeroLogger) Error(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {

	l.logger.
		Error().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Errorf(template string, args ...interface{}) {
	l.logger.
		Error().
		Msgf(template, args...)
}

func (l *zeroLogger) Fatal(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {

	l.logger.
		Fatal().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Fatalf(template string, args ...interface{}) {
	l.logger.
		Fatal().
		Msgf(template, args...)
}
