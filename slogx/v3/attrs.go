package slogx

import (
	"log/slog"
	"time"
)

var ErrorKey = "error"

func ErrorString(err error) Attr {
	return slog.String(ErrorKey, err.Error())
}

func ErrorValue(err error) slog.Value {
	return slog.StringValue(err.Error())
}

// TimeValue returns a Value for a time.Time.
// It discards the monotonic portion.
func TimeValue(v time.Time) any {
	return uint64(v.UnixNano())
}

// DurationValue returns a Value for a time.Duration.
func DurationValue(v time.Duration) uint64 {
	return uint64(v.Nanoseconds())
}

func String(key, v string) Attr {
	return slog.String(key, v)
}

func Int64(key string, v int64) Attr {
	return slog.Int64(key, v)
}

func Int(key string, v int) Attr {
	return slog.Int(key, v)
}

// Uint64 returns an Attr for a uint64.
func Uint64(key string, v uint64) Attr {
	return slog.Uint64(key, v)
}

func Float64(key string, v float64) Attr {
	return slog.Float64(key, v)
}

func Bool(key string, v bool) Attr {
	return slog.Bool(key, v)
}

func Time(key string, v time.Time) Attr {
	return slog.Time(key, v)
}

func Duration(key string, v time.Duration) Attr {
	return slog.Duration(key, v)
}

func Group(key string, args ...any) Attr {
	return slog.Group(key, args...)
}

func Any(key string, v any) Attr {
	return slog.Any(key, v)
}
