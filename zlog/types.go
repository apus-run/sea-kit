package zlog

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Panicf(format string, args ...any)

	Info(msg string, tags ...Field)
	Error(msg string, tags ...Field)
	Debug(msg string, tags ...Field)
	Warn(msg string, tags ...Field)
	Fatal(msg string, tags ...Field)
	Panic(msg string, tags ...Field)

	Slow(msg string, fields ...Field)
	Stack(msg string)
	Stat(msg string, fields ...Field)

	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)

	Close() error
	Sync() error
}

type Field struct {
	Key   string
	Value any
}
