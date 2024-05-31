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

	"github.com/apus-run/sea-kit/errorsx"
	"github.com/apus-run/sea-kit/slogx/v2/prettylog"
)

const LOGGER_KEY = "slogLogger"
const AttrErrorKey = "error"

type Logger struct {
	*slog.Logger
}

// NewLogger 只包装了 slog
func NewLogger(l *slog.Logger) *Logger {
	return &Logger{l}
}

// Log send log records with caller depth
func (l *Logger) Log(ctx context.Context, depth int, err error, level slog.Level, msg string, attrs ...any) {
	if !l.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	if err != nil {
		r.Add(AttrErrorKey, err)
	}
	r.Add(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}

var defaultHandler *Handler

// New 包装了 slog
func New(options ...Option) *slog.Logger {
	opts := Apply(options...)

	var handler slog.Handler
	handlerOptions := &slog.HandlerOptions{
		Level:       getLogLevel(opts.LogLevel),
		AddSource:   true,
		ReplaceAttr: ReplaceAttr,
	}

	switch f := opts.Format; f {
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	case FormatJSON:
		handler = slog.NewJSONHandler(opts.Writer, handlerOptions)
	case FormatPretty:
		handler = prettylog.NewHandler(&slog.HandlerOptions{
			Level:     getLogLevel(opts.LogLevel),
			AddSource: true,
		})
	default:
		handler = slog.NewJSONHandler(opts.Writer, handlerOptions)
	}

	if opts.LogGroup != "" {
		handler = handler.WithGroup(opts.LogGroup)
	}
	if len(opts.LogAttrs) > 0 {
		handler = handler.WithAttrs(opts.LogAttrs)
	}

	defaultHandler = NewHandler(handler).(*Handler)

	return slog.New(defaultHandler)
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
		return slog.LevelError
	}

	return *level
}

func ApplyHandlerOption(opt HandlerOption) {
	defaultHandler.Apply(opt)
}

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
				slog.Any("stack", errorsx.StackTrace(st)),
			))
		}
		return a
	}
	return a
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
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
