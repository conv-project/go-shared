package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// GetLogger returns singleton instance of logger.
func GetLogger() *zap.Logger {
	return instance
}

// InitLogger initializes logger with specified level.
func InitLogger(level string) error {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}

	once.Do(func() {
		config := zap.NewProductionEncoderConfig()
		config.TimeKey = "timestamp"
		config.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			zapcore.AddSync(os.Stdout),
			logLevel,
		)

		instance = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})

	return nil
}

// Info logs message at info level.
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Error logs message at error level.
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Debug logs message at debug level.
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Warn logs message at warn level.
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal logs message at fatal level and then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
