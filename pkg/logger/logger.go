package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config interface {
	GetLogLevel() string
	GetServerMode() string
}

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

type zapLogger struct {
	logger *zap.Logger
}

type noopLogger struct{}

func (n *noopLogger) Info(msg string, fields ...zap.Field)  {}
func (n *noopLogger) Error(msg string, fields ...zap.Field) {}
func (n *noopLogger) Warn(msg string, fields ...zap.Field)  {}
func (n *noopLogger) Debug(msg string, fields ...zap.Field) {}
func (n *noopLogger) Fatal(msg string, fields ...zap.Field) {}
func (n *noopLogger) With(fields ...zap.Field) Logger       { return n }
func (n *noopLogger) Sync() error                           { return nil }

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func NewLogger(env string, logLevel string) (Logger, error) {
	if strings.ToLower(logLevel) == "none" || strings.ToLower(logLevel) == "disabled" {
		return &noopLogger{}, nil
	}

	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if logLevel != "" {
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

func NewLoggerWithConfig(cfg Config) (Logger, error) {
	return NewLogger(cfg.GetServerMode(), cfg.GetLogLevel())
}

func NewProductionLogger() (Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	return NewLogger("production", logLevel)
}

func NewDevelopmentLogger() (Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "debug"
	}
	return NewLogger("development", logLevel)
}

func NewLoggerFromEnv() (Logger, error) {
	env := os.Getenv("GIN_MODE")
	if env == "" {
		env = "development"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "debug"
	}
	return NewLogger(env, logLevel)
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
