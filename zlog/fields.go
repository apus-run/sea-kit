package zlog

import (
	"time"
)

// TimeValue returns a Value for a time.Time.
// It discards the monotonic portion.
func TimeValue(v time.Time) any {
	return uint64(v.UnixNano())
}

// DurationValue returns a Value for a time.Duration.
func DurationValue(v time.Duration) uint64 {
	return uint64(v.Nanoseconds())
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func String(key, v string) Field {
	return Field{
		Key:   key,
		Value: v,
	}
}

func Uint64(key string, v uint64) Field {
	return Field{
		Key:   key,
		Value: v,
	}
}

func Int64(key string, v int64) Field {
	return Field{
		Key:   key,
		Value: v,
	}
}

func Float64(key string, v float64) Field {
	return Field{
		Key:   key,
		Value: v,
	}
}

func Bool(key string, b bool) Field {
	return Field{
		Key:   key,
		Value: b,
	}
}

func Int(key string, v int) Field {
	return Int64(key, int64(v))
}

func Time(key string, v time.Time) Field {
	return Field{Key: key, Value: TimeValue(v)}
}

func Duration(key string, v time.Duration) Field {
	return Field{Key: key, Value: DurationValue(v)}
}
