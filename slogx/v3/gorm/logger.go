package gorm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	log "github.com/apus-run/van/log/slog"
)

// Logger adapts to gorm logger
type Logger struct {
	*log.SlogLogger
	*Config
}

func NewLogger(options ...Option) *Logger {

	// 修改配置
	config := Apply(options...)
	l := log.NewLogger(log.WithOptions(func(lopts *log.Options) {
		lopts.LogLevel = config.LogLevel
		lopts.Format = config.Format
		lopts.Writer = config.Writer
		lopts.LogGroup = config.LogGroup
		lopts.LogAttrs = config.LogAttrs
	}))

	return &Logger{
		l,
		config,
	}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, i ...interface{}) {
	l.Log(ctx, 5, nil, slog.LevelInfo, msg, i...)
}

func (l *Logger) Warn(ctx context.Context, msg string, i ...interface{}) {
	l.Log(ctx, 5, nil, slog.LevelWarn, msg, i...)
}

func (l *Logger) Error(ctx context.Context, msg string, i ...interface{}) {
	l.Log(ctx, 5, nil, slog.LevelError, msg, i...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.Log(ctx, 5, err, slog.LevelError, "query error",
			slog.String("elapsed", elapsed.String()),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold:
		msg := fmt.Sprintf("slow threshold >= %v", l.SlowThreshold)
		l.Log(ctx, 5, nil, slog.LevelWarn, msg,
			slog.String("elapsed", elapsed.String()),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	case l.LogInfo:
		l.Log(ctx, 5, nil, slog.LevelInfo, "query info",
			slog.String("elapsed", elapsed.String()),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	}
}
