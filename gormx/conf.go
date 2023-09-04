package gormx

import (
	"fmt"

	"gorm.io/gorm"
)

// DriverType is the client driver
type DriverType int

// The Driver Type of native client
const (
	Unknown DriverType = iota
	MySQL
	PostgreSQL
	SQLite
	SQLServer
	ClickHouse
)

// DriverMap is the map of [driver, name]
var DriverMap = map[DriverType]string{
	MySQL:      "mysql",
	PostgreSQL: "postgres",
	SQLite:     "sqlite",
	SQLServer:  "sqlserver",
	ClickHouse: "clickhouse",
	Unknown:    "unknown",
}

// DriverTypeMap is the map of driver [name, driver]
var DriverTypeMap = ReverseMap(DriverMap)

// String convert the DriverType to string
func (d *DriverType) String() string {
	if val, ok := DriverMap[*d]; ok {
		return val
	}
	return DriverMap[Unknown]
}

// DriverType convert the string to DriverType
func (d *DriverType) DriverType(name string) DriverType {
	if val, ok := DriverTypeMap[name]; ok {
		*d = val
		return val
	}
	return Unknown
}

// DBOption 代表初始化的时候的选项
type DBOption func(*DBConfig) error

// DBConfig is the database configuration
type DBConfig struct {
	DriverType DriverType `json:"driver"`
	Dsn        string     `json:"dsn"`

	// 以下配置关于连接池的配置，如果不设置，则使用默认值
	ConnMaxOpen     int    `json:"conn_max_open"`      // 最大连接数；-1：不限；默认：20
	ConnMaxIdle     int    `json:"conn_max_idle"`      // 最大空闲连接数；-1：不限；默认：10
	ConnMaxLifetime string `json:"conn_max_lifetime"`  // 连接最大生命周期；-1：不限；默认：10分钟
	ConnMaxIdleTime string `json:"conn_max_idle_time"` // 空闲最大生命周期；-1：不限；默认：5分钟

	// 集成gorm的配置
	*gorm.Config
}

func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Dsn:             "",
		DriverType:      Unknown,
		ConnMaxOpen:     100,
		ConnMaxIdle:     10,
		ConnMaxLifetime: "300ms",
		ConnMaxIdleTime: "1ms",
		Config:          &gorm.Config{},
	}
}

func Apply(opts ...DBOption) *DBConfig {
	config := DefaultDBConfig()
	for _, opt := range opts {
		err := opt(config)
		if err != nil {
			return nil
		}
	}
	return config
}

// WithGormConfig 表示自行配置Gorm的配置信息
func WithGormConfig(f func(options *DBConfig)) DBOption {
	return func(config *DBConfig) error {
		f(config)
		return nil
	}
}

// Check do the configuration check
func (d *DBConfig) Check() error {
	if d.DriverType == Unknown {
		return fmt.Errorf("unknown driver")
	}
	if d.ConnMaxOpen < 0 {
		return fmt.Errorf("conn_max_open must be greater than 0")
	}

	if d.ConnMaxIdle < 0 {
		return fmt.Errorf("conn_max_idle must be greater than 0")
	}

	return nil
}

// ReverseMap just reverse the map from [key, value] to [value, key]
func ReverseMap[K comparable, V comparable](m map[K]V) map[V]K {
	n := make(map[V]K, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
