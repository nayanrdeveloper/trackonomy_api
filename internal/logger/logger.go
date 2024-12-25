package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger initializes the global zap logger.
// In production, you might configure it differently than in dev.
func InitLogger(isProd bool) error {
	var cfg zap.Config
	if isProd {
		// Production Config: JSON format, Info level, etc.
		cfg = zap.NewProductionConfig()
		// You can also modify the default config here, e.g. changing log level:
		// cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		// Development Config: human-readable format, Debug level, colored output
		cfg = zap.NewDevelopmentConfig()

		// Example overrides in development:
		// 1) Only show stack traces on Error level or higher (instead of Warn)
		cfg.EncoderConfig.StacktraceKey = "stacktrace"   // rename key or set to ""
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel) // set minimum log level
		cfg.DisableStacktrace = false                    // by default false in dev
		// 2) Control how levels are displayed (capital with color)
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Build logger
	l, err := cfg.Build()
	if err != nil {
		return err
	}
	log = l
	return nil
}

// Sync flushes any buffered log entries. Call this before your app exits.
func Sync() {
	_ = log.Sync()
}

// For convenience, expose zap’s methods in a structured way:

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// If you need to add fields frequently, you can create
// convenience wrappers or use zap.With(...) to attach fields persistently.
