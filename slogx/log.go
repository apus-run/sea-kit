package slogx

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
)

const LOGGER_KEY = "slogLogger"

var factory = &LoggerFactory{
	loggers: make(map[string]*SlogLogger),
}

type LoggerFactory struct {
	mu      sync.Mutex
	loggers map[string]*SlogLogger
}

type SlogLogger struct {
	logger *slog.Logger
}

func NewLogger(opts ...Option) *SlogLogger {
	options := Apply(opts...)

	factory.mu.Lock()
	if logger, ok := factory.loggers[options.logFilename]; ok {
		factory.mu.Unlock()
		return logger
	}
	defer factory.mu.Unlock()

	// 日志文件切割归档
	writerSyncer := getLogWriter(options)

	// 日志级别
	level := getLogLevel(options.logLevel)

	var handler slog.Handler
	if len(options.logFilename) == 0 && options.encoding == "console" {
		handler = textHandler(os.Stdout, level)
	} else {
		handler = jsonHandler(writerSyncer, level)
	}

	l := slog.New(handler)
	logger := &SlogLogger{l}
	factory.loggers[options.logFilename] = logger

	return logger
}

func getLogWriter(opts *Options) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   opts.logFilename,
		MaxSize:    opts.maxSize, // megabytes
		MaxBackups: opts.maxBackups,
		MaxAge:     opts.maxAge, //days
		Compress:   opts.compress,
	}
}

func jsonHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
}

func textHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewTextHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
}

func getLogLevel(logLevel string) slog.Level {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return slog.LevelError
	}

	return *level
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
