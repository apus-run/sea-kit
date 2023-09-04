package sqlx

import (
	"time"

	"github.com/apus-run/sea-kit/log"
)

// Option is database option
type Option func(*options)

type options struct {
	dsn             string        // 数据库连接地址
	maxOpenConns    int           // default: 100
	maxIdleConns    int           // default: 10
	connMaxLifetime time.Duration // default: 300s
	logging         bool          // default: "false"

	logger log.Logger
}

// DefaultOptions .
func DefaultOptions() *options {
	return &options{
		dsn:             "",
		maxOpenConns:    100,
		maxIdleConns:    10,
		connMaxLifetime: 300 * time.Minute,
		logging:         false,

		logger: log.DefaultLogger,
	}
}

// WithDSN .
func WithDSN(dsn string) Option {
	return func(o *options) {
		o.dsn = dsn
	}
}

// WithMaxOpenConns .
func WithMaxOpenConns(moc int) Option {
	return func(o *options) {
		o.maxOpenConns = moc
	}
}

// WithMaxIdleConns .
func WithMaxIdleConns(mic int) Option {
	return func(o *options) {
		o.maxIdleConns = mic
	}
}

// WithConnMaxLifetime .
func WithConnMaxLifetime(cml time.Duration) Option {
	return func(o *options) {
		o.connMaxLifetime = cml
	}
}

// WithLogging .
func WithLogging(logging bool) Option {
	return func(o *options) {
		o.logging = logging
	}
}

// Logger .
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}
