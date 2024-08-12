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
	*zap.Logger
}

// NewZapLogger 只包装了 zap
func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{
		l,
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
			zap.New(core, zap.Development(), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel)),
		}
	}
	return &ZapLogger{
		zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel)),
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
		// 创建一个默认的 encoder 配置
		encoderConfig := zap.NewProductionEncoderConfig()
		// 自定义 MessageKey 为 message，message 语义更明确
		encoderConfig.MessageKey = "message"
		// 自定义 TimeKey 为 timestamp，timestamp 语义更明确
		encoderConfig.TimeKey = "timestamp"
		// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// 指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000` 格式，更易读
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		// 指定 time.Duration 序列化函数，将 time.Duration 序列化为经过的毫秒数的浮点数
		// 毫秒数比默认的秒数更精确
		encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendFloat64(float64(d) / float64(time.Millisecond))
		}
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

func (l *ZapLogger) Info(msg string, tags ...Field) {
	l.Logger.Info(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Error(msg string, tags ...Field) {
	l.Logger.Error(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Debug(msg string, tags ...Field) {
	l.Logger.Debug(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Warn(msg string, tags ...Field) {
	l.Logger.Warn(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Fatal(msg string, tags ...Field) {
	l.Logger.Fatal(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Panic(msg string, tags ...Field) {
	l.Logger.Panic(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Slow(msg string, tags ...Field) {
	l.Logger.Warn(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Stack(msg string) {
	l.Logger.Error(fmt.Sprint(msg), zap.Stack("stack"))
}

func (l *ZapLogger) Stat(msg string, tags ...Field) {
	l.Logger.Info(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	l.Logger.Sugar().Debugf(format, args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	l.Logger.Sugar().Infof(format, args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	l.Logger.Sugar().Warnf(format, args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	l.Logger.Sugar().Errorf(format, args...)
}

func (l *ZapLogger) Fatalf(format string, args ...any) {
	l.Logger.Sugar().Fatalf(format, args...)
}

func (l *ZapLogger) Panicf(format string, args ...any) {
	l.Logger.Sugar().Panicf(format, args...)
}

func (l *ZapLogger) Print(args ...any) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *ZapLogger) Printf(format string, args ...any) {
	l.Logger.Sugar().Infof(format, args...)
}

func (l *ZapLogger) Println(args ...any) {
	l.Logger.Info(fmt.Sprintln(args...))
}

func (l *ZapLogger) Close() error {
	return l.Logger.Sync()
}

func (l *ZapLogger) Sync() error {
	return l.Logger.Sync()
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *ZapLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return l
	}

	lc := l.clone()
	lc.Logger = lc.Logger.With(l.toZapFields(fields)...)
	return lc
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents zap from always
// reporting the wrapper code as the caller.
func (l *ZapLogger) AddCallerSkip(skip int) Logger {
	lc := l.clone()
	lc.Logger = lc.Logger.WithOptions(zap.AddCallerSkip(skip))
	return lc
}

// clone 深度拷贝 zapLogger.
func (l *ZapLogger) clone() *ZapLogger {
	copied := *l
	return &copied
}

func (l *ZapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, arg := range fields {
		zapFields = append(zapFields, zap.Any(arg.Key, arg.String))
	}
	return zapFields
}
