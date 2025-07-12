// https://tonybai.com/2023/09/01/slog-a-new-choice-for-logging-in-go/

package slogx

import (
	"context"
	"fmt"
	"go/build"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const LOGGER_KEY = "slogLogger"

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelPanic = slog.Level(10)
	LevelFatal = slog.Level(12)
	LevelTrace = slog.Level(-8)
)

var LevelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
	LevelError: "ERROR",
	LevelWarn:  "WARN",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
	LevelPanic: "PANIC",
	LevelTrace: "TRACE",
}

type SlogLogger struct {
	*slog.Logger
}

var defaultHandler *Handler

// NewSlogLogger 只包装了 slog
func NewSlogLogger(l *slog.Logger) Logger {
	if l == nil {
		l = slog.Default()
	}
	return &SlogLogger{l}
}

func NewLogger(opts ...Option) *SlogLogger {
	options := Apply(opts...)

	handler := createHandler(options)
	logger := slog.New(handler)

	return &SlogLogger{
		Logger: logger,
	}
}

func createHandler(options *Options) slog.Handler {
	var handler slog.Handler
	// 日志级别
	level := getLogLevel(options.LogLevel)

	handlerOptions := &slog.HandlerOptions{
		Level:       level,
		AddSource:   true,
		ReplaceAttr: ReplaceAttr,
	}

	switch f := options.Format; f {
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	case FormatJSON:
		handler = slog.NewJSONHandler(options.Writer, handlerOptions)
	default:
		handler = slog.NewJSONHandler(options.Writer, handlerOptions)
	}

	if options.LogGroup != "" {
		handler = handler.WithGroup(options.LogGroup)
	}
	if len(options.LogAttrs) > 0 {
		handler = handler.WithAttrs(options.LogAttrs)
	}
	defaultHandler = NewHandler(handler).(*Handler)

	return defaultHandler
}

// NewNop returns a no-op logger
func NewNop() *slog.Logger {
	nopLevel := slog.Level(-99)
	ops := &slog.HandlerOptions{
		Level: nopLevel,
	}
	handler := slog.NewTextHandler(io.Discard, ops)
	return slog.New(handler)
}

// NewWithHandler build *slog.Logger with slog Handler
func NewWithHandler(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

func getLogLevel(logLevel string) slog.Level {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return LevelError
	}

	return *level
}

func ApplyHandlerOption(opt HandlerOption) {
	defaultHandler.Apply(opt)
}

const AttrErrorKey = "error"

// ReplaceAttr handle log key-value pair
func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		return slog.String(a.Key, a.Value.Time().Format(time.RFC3339))
	case slog.LevelKey:
		return slog.String(a.Key, strings.ToLower(a.Value.String()))
	case slog.SourceKey:
		if v, ok := a.Value.Any().(*slog.Source); ok {
			a.Value = slog.StringValue(fmt.Sprintf("%s:%d", getBriefSource(v.File), v.Line))
		}
		return a
	case AttrErrorKey:
		v, ok := a.Value.Any().(interface {
			StackTrace() errors.StackTrace
		})
		if ok {
			st := v.StackTrace()
			return slog.Any(a.Key, slog.GroupValue(
				slog.String("msg", a.Value.String()),
				slog.Any("stack", errors.StackTrace(st)),
			))
		}
		return a
	}
	return a
}

func projectPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	}
	return ""
}

func getBriefSource(source string) string {
	gp := filepath.ToSlash(build.Default.GOPATH)
	if strings.HasPrefix(source, gp) {
		return strings.TrimPrefix(source, gp+"/pkg/mod/")
	}
	pp := filepath.ToSlash(projectPath())
	return strings.TrimPrefix(source, pp+"/")
}

// Log send log records with caller depth
func (l *SlogLogger) Log(ctx context.Context, depth int, err error, level slog.Level, msg string, attrs ...any) {
	if !l.Enabled(ctx, level) {
		return
	}

	// 记录日志
	l.Logger.Log(ctx, level, msg, attrs...)
}

func (l *SlogLogger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...Attr) {
	if !l.Enabled(ctx, level) {
		return
	}
	l.Logger.LogAttrs(ctx, level, msg, attrs...)
}

// 实现 Logger 接口的方法
func (l *SlogLogger) Info(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelInfo, msg, attrs...)
}

func (l *SlogLogger) Infof(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelInfo, sprintf(msg, args...))
}

func (l *SlogLogger) InfoContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelInfo, msg, attrs...)
}

func (l *SlogLogger) Error(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelError, msg, attrs...)
}

func (l *SlogLogger) Errorf(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelError, sprintf(msg, args...))
}

func (l *SlogLogger) ErrorContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelError, msg, attrs...)
}

func (l *SlogLogger) Debug(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelDebug, msg, attrs...)
}

func (l *SlogLogger) Debugf(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelDebug, sprintf(msg, args...))
}

func (l *SlogLogger) DebugContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelDebug, msg, attrs...)
}

func (l *SlogLogger) Warn(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelWarn, msg, attrs...)
}

func (l *SlogLogger) Warnf(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelWarn, sprintf(msg, args...))
}

func (l *SlogLogger) WarnContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelWarn, msg, attrs...)
}

func (l *SlogLogger) Fatal(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelFatal, msg, attrs...)
	os.Exit(1)
}

func (l *SlogLogger) FatalContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelFatal, msg, attrs...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalf(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelFatal, sprintf(msg, args...))
	os.Exit(1)
}

func (l *SlogLogger) Panic(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelPanic, msg, attrs...)
	panic(SprintfWithAttrs(msg, attrs...))
}

func (l *SlogLogger) PanicContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelPanic, msg, attrs...)
	panic(SprintfWithAttrs(msg, attrs...))
}

func (l *SlogLogger) Panicf(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelPanic, sprintf(msg, args...))
	panic(fmt.Sprintf(msg, args...))
}

func (l *SlogLogger) Trace(msg string, attrs ...Attr) {
	l.Logger.LogAttrs(context.Background(), LevelTrace, msg, attrs...)
}
func (l *SlogLogger) Tracef(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, sprintf(msg, args...))
}
func (l *SlogLogger) TraceContext(ctx context.Context, msg string, attrs ...Attr) {
	l.Logger.LogAttrs(ctx, LevelTrace, msg, attrs...)
}

func SprintfWithAttrs(format string, attrs ...Attr) string {
	// 处理 Attr 类型并将其转换为字符串
	attrStr := ""
	for _, attr := range attrs {
		attrStr += fmt.Sprintf("%v ", attr) // 根据需要格式化 Attr
	}
	return fmt.Sprintf(format, attrStr)
}

func (l *SlogLogger) With(args ...any) *SlogLogger {
	lc := l.clone()
	lc.Logger = lc.Logger.With(args...)
	return lc
}

func (l *SlogLogger) New(args ...any) *SlogLogger {
	return l.With(args...)
}

func (l *SlogLogger) WithGroup(name string) *SlogLogger {
	lc := l.clone()
	lc.Logger = lc.Logger.WithGroup(name)
	return lc
}

// Enabled reports whether l emits log records at the given context and level.
func (l *SlogLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.Logger.Enabled(ctx, level)
}

// clone 深度拷贝 SlogLogger.
func (l *SlogLogger) clone() *SlogLogger {
	copied := *l
	return &copied
}
