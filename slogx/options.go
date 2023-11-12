package slogx

// Option is config option.
type Option func(*Options)

type Options struct {
	// logger options
	LogLevel string // debug, info, warn, error
	Encoding string // console or json

	// lumberjack options
	LogFilename string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		LogLevel: "info",
		Encoding: "console",

		LogFilename: "logs.log",
		MaxSize:     500, // megabytes
		MaxBackups:  3,
		MaxAge:      28, //days
		Compress:    true,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithLogLevel 日志级别
func WithLogLevel(level string) Option {
	return func(o *Options) {
		o.LogLevel = level
	}
}

// WithEncoding 日志编码
func WithEncoding(encoding string) Option {
	return func(o *Options) {
		o.Encoding = encoding
	}
}

// WithFilename 日志文件
func WithFilename(filename string) Option {
	return func(o *Options) {
		o.LogFilename = filename
	}
}

// WithMaxSize 日志文件大小
func WithMaxSize(maxSize int) Option {
	return func(o *Options) {
		o.MaxSize = maxSize
	}
}

// WithMaxBackups 日志文件最大备份数
func WithMaxBackups(maxBackups int) Option {
	return func(o *Options) {
		o.MaxBackups = maxBackups
	}
}

// WithMaxAge 日志文件最大保存时间
func WithMaxAge(maxAge int) Option {
	return func(o *Options) {
		o.MaxAge = maxAge
	}
}

// WithCompress 日志文件是否压缩
func WithCompress(compress bool) Option {
	return func(o *Options) {
		o.Compress = compress
	}
}
