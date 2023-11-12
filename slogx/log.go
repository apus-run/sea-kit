// https://tonybai.com/2023/09/01/slog-a-new-choice-for-logging-in-go/

package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/natefinch/lumberjack"
)

const LOGGER_KEY = "slogLogger"

var factory = &LoggerFactory{
	loggers: make(map[string]*slog.Logger),
}

type LoggerFactory struct {
	mu      sync.Mutex
	loggers map[string]*slog.Logger
}

var defaultHandler *Handler

func NewLogger(opts ...Option) *slog.Logger {
	options := Apply(opts...)

	factory.mu.Lock()
	if logger, ok := factory.loggers[options.LogFilename]; ok {
		factory.mu.Unlock()
		return logger
	}
	defer factory.mu.Unlock()

	// 日志文件切割归档
	writerSyncer := getLogWriter(options)

	// 日志级别
	level := getLogLevel(options.LogLevel)

	var handler slog.Handler
	handlerOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(fmt.Sprintf("%s:%d", src.File, src.Line))
				}
			}
			return a
		},
	}
	if len(options.LogFilename) == 0 && strings.ToLower(options.Encoding) == "console" {
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	} else {
		handler = slog.NewJSONHandler(writerSyncer, handlerOptions)
	}

	defaultHandler = NewHandler(handler).(*Handler)
	logger := slog.New(defaultHandler)
	// 此处设置默认日志, 最好手动设置
	// slog.SetDefault(l)

	factory.loggers[options.LogFilename] = logger

	logger.Info("the log module has been initialized successfully.", slog.Any("options", options))

	return logger
}

func getLogWriter(opts *Options) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   opts.LogFilename,
		MaxSize:    opts.MaxSize, // megabytes
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge, //days
		Compress:   opts.Compress,
	}
}

func getLogLevel(logLevel string) slog.Level {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return slog.LevelError
	}

	return *level
}

func ApplyHandlerOption(opt HandlerOption) {
	defaultHandler.Apply(opt)
}
