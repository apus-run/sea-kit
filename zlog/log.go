package zlog

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOGGER_KEY = "zapLogger"

type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger 只包装了 zap
func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{
		logger: l,
	}
}

// NewLogger 包装了 zap 和日志文件切割归档
func NewLogger(opts ...Option) *ZapLogger {
	options := Apply(opts...)

	// 日志文件切割归档
	// writerSyncer := getLogWriter(opts...)
	writerSyncer := getLogConsoleWriter(options)

	// 编码器配置
	encoder := getEncoder(options.Encoding)

	// 日志级别
	level := getLogLevel(options.LogLevel)

	core := zapcore.NewCore(encoder, writerSyncer, level)
	if options.Mode != "prod" {
		return &ZapLogger{
			zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
		}
	}
	return &ZapLogger{
		zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
	}
}

func getEncoder(encoding string) zapcore.Encoder {
	if encoding == "console" {
		// NewConsoleEncoder 打印更符合人们观察的方式
		return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 在日志文件中使用大写字母记录日志级别
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

func getLogWriter(opts *Options) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   opts.LogFilename,
		MaxSize:    opts.MaxSize, // megabytes
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge, //days
		Compress:   opts.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLogConsoleWriter(opts *Options) zapcore.WriteSyncer {
	// 日志文件切割归档
	lumberJackLogger := &lumberjack.Logger{
		Filename:   opts.LogFilename,
		MaxSize:    opts.MaxSize, // megabytes
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge, //days
		Compress:   opts.Compress,
	}

	// 打印到控制台和文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}

func getLogLevel(logLevel string) zapcore.Level {
	level := new(zapcore.Level)
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return zap.ErrorLevel
	}

	return *level
}

func (l *ZapLogger) With(args ...Field) Logger {
	z := l.logger.With(l.toZapFields(args)...)
	return NewZapLogger(z)
}

func (l *ZapLogger) Info(msg string, tags ...Field) {
	l.logger.Info(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Error(msg string, tags ...Field) {
	l.logger.Error(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Debug(msg string, tags ...Field) {
	l.logger.Debug(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Warn(msg string, tags ...Field) {
	l.logger.Warn(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Fatal(msg string, tags ...Field) {
	l.logger.Fatal(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Panic(msg string, tags ...Field) {
	l.logger.Panic(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Slow(msg string, tags ...Field) {
	l.logger.Warn(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Stack(msg string) {
	l.logger.Error(fmt.Sprint(msg), zap.Stack("stack"))
}

func (l *ZapLogger) Stat(msg string, tags ...Field) {
	l.logger.Info(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Debug(msg, zap.Any("args", args))
}

func (l *ZapLogger) Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Info(msg, zap.Any("args", args))
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Warn(msg, zap.Any("args", args))
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Error(msg, zap.Any("args", args))
}

func (l *ZapLogger) Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Fatal(msg, zap.Any("args", args))
}

func (l *ZapLogger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Panic(msg, zap.Any("args", args))
}

func (l *ZapLogger) Print(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *ZapLogger) Printf(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Println(args ...any) {
	l.logger.Info(fmt.Sprintln(args...))
}

func (l *ZapLogger) Close() error {
	return l.logger.Sync()
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *ZapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, arg := range fields {
		zapFields = append(zapFields, zap.Any(arg.Key, arg.Value))
	}
	return zapFields
}
