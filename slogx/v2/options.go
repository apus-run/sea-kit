package slogx

import (
	"io"
	"log/slog"
	"os"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	LogLevel string      // debug, info, warn, error
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
func WithLogAttrs(attrs []slog.Attr) Option {
	return func(o *Options) {
		o.LogAttrs = attrs
	}
}

// WithFormat set log format
func WithFormat(format Format) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// WithWriter set log writer
func WithWriter(writer io.Writer) Option {
	return func(o *Options) {
		o.Writer = writer
	}
}

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)
