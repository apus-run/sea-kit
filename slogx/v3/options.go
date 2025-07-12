package slogx

import (
	"io"
	"log/slog"
	"os"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	// logger options
	LogLevel string      // debug, info, warn, error,debug, panic
	Format   Format      // text or json
	Writer   io.Writer   // 日志输出
	LogGroup string      // slog group
	LogAttrs []slog.Attr // 日志属性
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		LogLevel: "info",
		Format:   FormatText,
		Writer:   os.Stdout,
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

// WithLogGroup 日志分组
func WithLogGroup(group string) Option {
	return func(o *Options) {
		o.LogGroup = group
	}
}

// WithLogAttrs 日志属性
func WithLogAttrs(attrs []Attr) Option {
	return func(o *Options) {
		o.LogAttrs = attrs
	}
}

// WithFormat 日志格式
func WithFormat(format Format) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// WithWriter 日志输出
func WithWriter(writer io.Writer) Option {
	return func(o *Options) {
		o.Writer = writer
	}
}

// WithOptions 设置所有配置
func WithOptions(fn func(options *Options)) Option {
	return func(options *Options) {
		fn(options)
	}
}
