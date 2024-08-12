package gorm

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/apus-run/sea-kit/zlog/v2"
	gormlogger "gorm.io/gorm/logger"
)

// Logger adapts to zlog logger
type Logger struct {
	*zlog.ZapLogger
	*Config
}

func NewLogger(options ...Option) *Logger {
	// 修改配置
	config := Apply(options...)
	z := zlog.NewLogger(zlog.WithOptions(func(zopts *zlog.Options) {
		zopts.Encoding = config.Encoding
		zopts.LogFilename = config.LogFilename
		zopts.MaxSize = config.MaxSize
		zopts.MaxBackups = config.MaxBackups
		zopts.MaxAge = config.MaxAge
		zopts.Compress = config.Compress
		zopts.Mode = config.Mode
		zopts.LogLevel = config.LogLevel
	}))
	return &Logger{
		z,
		config,
	}
}

var (
	infoStr       = "%s[info] "
	warnStr       = "%s[warn] "
	errStr        = "%s[error] "
	traceStr      = "[%s][%.3fms] [rows:%v] %s"
	traceWarnStr  = "%s %s[%.3fms] [rows:%v] %s"
	traceErrStr   = "%s %s[%.3fms] [rows:%v] %s"
	slowThreshold = 200 * time.Millisecond
)

var levelM = map[string]gormlogger.LogLevel{
	"panic": gormlogger.Silent,
	"error": gormlogger.Error,
	"warn":  gormlogger.Warn,
	"info":  gormlogger.Info,
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, keyvals ...any) {
	l.Infof(infoStr+msg, append([]any{fileWithLineNum()}, keyvals...)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, keyvals ...any) {
	l.Warnf(warnStr+msg, append([]any{fileWithLineNum()}, keyvals...)...)
}

func (l *Logger) Error(ctx context.Context, msg string, keyvals ...any) {
	l.Errorf(errStr+msg, append([]any{fileWithLineNum()}, keyvals...)...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if levelM[l.LogLevel] <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && levelM[l.LogLevel] >= gormlogger.Error:
		sql, rows := fc()
		if rows == -1 {
			l.Errorf(traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Errorf(traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > slowThreshold && levelM[l.LogLevel] >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold)
		if rows == -1 {
			l.Warnf(traceWarnStr, fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Warnf(traceWarnStr, fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case levelM[l.LogLevel] >= gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Infof(traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Infof(traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func fileWithLineNum() string {
	for i := 4; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)

		// if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
		if ok && !strings.HasSuffix(file, "_test.go") {
			dir, f := filepath.Split(file)

			return filepath.Join(filepath.Base(dir), f) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
