package gormx

import (
	"gorm.io/gorm"
	"log/slog"
	"time"
)

// Handler ...
type Handler func(*gorm.DB)

// Processor ...
type Processor interface {
	Get(name string) func(*gorm.DB)
	Replace(name string, handler func(*gorm.DB)) error
}

// Interceptor ...
type Interceptor func(string) func(next Handler) Handler

func debugInterceptor() func(Handler) Handler {
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			beg := time.Now()
			next(db)
			cost := time.Since(beg)
			if db.Error != nil {
				slog.Debug("error", db.Error, "cost", cost, "sql", logSQL(db, false))
			} else {
				slog.Debug("", "cost", cost, "sql", logSQL(db, true))
			}
		}
	}
}

func metricInterceptor() func(next Handler) Handler {
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			beg := time.Now()
			next(db)
			cost := time.Since(beg)
			slog.Debug("metric", "cost", cost, "sql", logSQL(db, true))

			// 写一些东西
		}
	}
}

func traceInterceptor() func(next Handler) Handler {
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			beg := time.Now()
			next(db)
			cost := time.Since(beg)
			slog.Debug("trace", "cost", cost, "sql", logSQL(db, true))

			// 写一些东西
		}
	}
}

func logSQL(db *gorm.DB, enableDetailSQL bool) string {
	if enableDetailSQL {
		return db.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	}
	return db.Statement.SQL.String()
}
