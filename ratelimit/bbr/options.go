package bbr

import (
	"time"
)

// options of bbr limiter.
type options struct {
	// WindowSize defines time duration per window
	Window time.Duration
	// BucketNum defines bucket number for each window
	Bucket int
	// CPUThreshold
	CPUThreshold int64
	// CPUQuota
	CPUQuota float64
}

// Option function for bbr limiter
type Option func(*options)

// WithWindow with window size.
func WithWindow(d time.Duration) Option {
	return func(o *options) {
		o.Window = d
	}
}

// WithBucket with bucket ize.
func WithBucket(b int) Option {
	return func(o *options) {
		o.Bucket = b
	}
}

// WithCPUThreshold with cpu threshold;
func WithCPUThreshold(threshold int64) Option {
	return func(o *options) {
		o.CPUThreshold = threshold
	}
}

// WithCPUQuota with real cpu quota(if it can not collect from process correct);
func WithCPUQuota(quota float64) Option {
	return func(o *options) {
		o.CPUQuota = quota
	}
}

// DefaultOptions .
func DefaultOptions() options {
	return options{
		Window:       time.Second * 10,
		Bucket:       100,
		CPUThreshold: 800,
	}
}

func Apply(opts ...Option) options {
	def := DefaultOptions()
	for _, apply := range opts {
		apply(&def)
	}
	return def
}
