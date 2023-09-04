package zlog

// Option is config option.
type Option func(*Options)

type Options struct {
	// logger options
	mode     string // dev or prod
	logLevel string // debug, info, warn, error, panic, panic, fatal
	encoding string // console or json

	// lumberjack options
	logFilename string
	maxSize     int
	maxBackups  int
	maxAge      int
	compress    bool
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		mode:     "dev",
		logLevel: "info",
		encoding: "console",

		logFilename: "logs.log",
		maxSize:     500, // megabytes
		maxBackups:  3,
		maxAge:      28, //days
		compress:    true,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithMode 运行模式
func WithMode(mode string) Option {
	return func(o *Options) {
		o.mode = mode
	}
}

// WithLogLevel 日志级别
func WithLogLevel(level string) Option {
	return func(o *Options) {
		o.logLevel = level
	}
}

// WithEncoding 日志编码
func WithEncoding(encoding string) Option {
	return func(o *Options) {
		o.encoding = encoding
	}
}

// WithFilename 日志文件
func WithFilename(filename string) Option {
	return func(o *Options) {
		o.logFilename = filename
	}
}

// WithMaxSize 日志文件大小
func WithMaxSize(maxSize int) Option {
	return func(o *Options) {
		o.maxSize = maxSize
	}
}

// WithMaxBackups 日志文件最大备份数
func WithMaxBackups(maxBackups int) Option {
	return func(o *Options) {
		o.maxBackups = maxBackups
	}
}

// WithMaxAge 日志文件最大保存时间
func WithMaxAge(maxAge int) Option {
	return func(o *Options) {
		o.maxAge = maxAge
	}
}

// WithCompress 日志文件是否压缩
func WithCompress(compress bool) Option {
	return func(o *Options) {
		o.compress = compress
	}
}
