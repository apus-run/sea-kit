package log

import (
	"log/slog"
	"time"
)

type Option func(l *Logger)

// WithLogger defines a custom logger to use
func WithLogger(log *slog.Logger) Option {
	return func(l *Logger) {
		l.slogger = log
	}
}

// WithSourceField defines the field to set the file name and line number of the current file
func WithSourceField(field string) Option {
	return func(l *Logger) {
		l.sourceField = field
	}
}

// WithErrorField defines the field to set the error
func WithErrorField(field string) Option {
	return func(l *Logger) {
		l.errorField = field
	}
}

// WithSlowThreshold defines the threshold above which a sql query is considered slow
func WithSlowThreshold(threshold time.Duration) Option {
	return func(l *Logger) {
		l.slowThreshold = threshold
	}
}

// WithTraceAll enables mode which logs all SQL messages.
func WithTraceAll() Option {
	return func(l *Logger) {
		l.traceAll = true
	}
}

// SetLogLevel sets a new slog.Level for a LogType.
func SetLogLevel(key LogType, level slog.Level) Option {
	return func(l *Logger) {
		l.logLevel[key] = level
	}
}

// WithRecordNotFoundError allows the slogger to log gorm.ErrRecordNotFound errors
func WithRecordNotFoundError() Option {
	return func(l *Logger) {
		l.ignoreRecordNotFoundError = false
	}
}

// WithIgnoreTrace disables the tracing of SQL queries by the slogger
func WithIgnoreTrace() Option {
	return func(l *Logger) {
		l.ignoreTrace = true
	}
}
