package sqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"runtime"
	"sync"

	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/xo/dburl"
	"golang.org/x/sync/singleflight"

	"github.com/apus-run/sea-kit/log"

	"github.com/apus-run/sea-kit/sqlx/hooks"
)

var (
	sfg singleflight.Group
	rwl sync.RWMutex

	dbs = map[string]*DB{}
)

type Database interface {
	Get(ctx context.Context, name string) *DB
}

type database struct {
	opts *options
	log  *log.Helper

	db *DB
}

func New(opts ...Option) Database {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &database{
		opts: options,
		log:  log.NewHelper(options.logger),
	}
}

// Get DB 获取数据库实例, db := db.Get(ctx, "foo")
func (data *database) Get(ctx context.Context, name string) *DB {
	rwl.RLock()

	if db, ok := dbs[name]; ok {
		rwl.RUnlock()
		return db
	}

	rwl.RUnlock()

	v, _, _ := sfg.Do(name, func() (interface{}, error) {
		u, err := dburl.Parse(data.opts.dsn)
		if err != nil {
			return nil, fmt.Errorf("unable to parse database URL: %v", err)
		}

		// 设置用户名和密码
		// u.User = url.UserPassword(db.opts.username, db.opts.password)

		sdb := sqlx.MustOpen(u.Driver, u.DSN)

		// Mapper function for SQL name mapping, snake_case table names
		sdb.MapperFunc(strcase.ToSnake)

		db := &DB{sdb}

		rwl.Lock()
		defer rwl.Unlock()
		dbs[name] = db

		return db, nil
	})

	return v.(*DB)
}

// GetWithHooks .
func (data *database) GetWithHooks(ctx context.Context, name string) *DB {
	rwl.RLock()

	if db, ok := dbs[name]; ok {
		rwl.RUnlock()
		return db
	}

	rwl.RUnlock()

	v, _, _ := sfg.Do(name, func() (interface{}, error) {
		var driverName string
		var d driver.Driver

		u, err := dburl.Parse(data.opts.dsn)
		if err != nil {
			return nil, fmt.Errorf("unable to parse database URL: %v", err)
		}

		driverName = fmt.Sprintf("%s:%s", u.Driver, name)
		hs := hooks.CombineHooks(
			&hooks.LogHooks{
				Log: data.log,
			},
			&hooks.MetricHooks{},
			&hooks.TracingHooks{},
		)

		switch u.Driver {
		case "sqlite":
			d = sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, hs)
		case "mysql":
			d = sqlhooks.Wrap(&mysql.MySQLDriver{}, hs)
		case "postgres":
			d = sqlhooks.Wrap(stdlib.GetDefaultDriver(), hs)
		default:
			d = sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, hs)
		}
		// 设置用户名和密码
		// u.User = url.UserPassword(db.opts.username, db.opts.password)

		sql.Register(driverName, d)
		sdb := sqlx.MustOpen(driverName, u.DSN)

		// Mapper function for SQL name mapping, snake_case table names
		sdb.MapperFunc(strcase.ToSnake)

		db := &DB{sdb}

		rwl.Lock()
		defer rwl.Unlock()
		dbs[name] = db

		return db, nil
	})

	return v.(*DB)
}

// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
func (data *database) getMaxOpenConnection() int {
	limit := data.opts.maxOpenConns

	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}

	if limit > 1024 {
		limit = 1024
	}

	return limit
}

// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
func (data *database) getMaxIdleConnection() int {
	limit := data.opts.maxIdleConns

	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}

	if limit > data.getMaxOpenConnection() {
		limit = data.getMaxOpenConnection()
	}

	return limit
}

// isConnect 是否连接成功
func (data *database) isConnect() bool {
	if err := data.db.Ping(); err == nil {
		return true
	}
	return false
}

func (data *database) GetClient() *DB {
	return data.db
}

// Close 关闭数据库连接，如果未指定名称，则关闭全部
func (data *database) Close(ctx context.Context, name ...string) {
	if len(name) == 0 {
		for key, db := range dbs {
			if err := db.Close(); err != nil {
				log.Error(fmt.Sprintf("db.%s close error", key), zap.Error(err))
			}
		}

		return
	}

	for _, key := range name {
		if db, ok := dbs[key]; ok {
			if err := db.Close(); err != nil {
				log.Error(fmt.Sprintf("db.%s close error", key), zap.Error(err))
			}
		}
	}
}
