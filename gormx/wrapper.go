package gormx

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/apus-run/sea-kit/log"
)

// DuplicateEntryErrCode 字段数据重复出现冲突, 状态码
const DuplicateEntryErrCode uint16 = 1062

var _ Wrapper = (*Helper)(nil)

type TableColumn struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}

type Wrapper interface {
	// GetDB 获取某个db
	GetDB(option ...DBOption) (*gorm.DB, error)

	// CanConnect 是否可以连接
	CanConnect(ctx context.Context, db *gorm.DB) (bool, error)

	// GetTables Table 相关
	GetTables(ctx context.Context, db *gorm.DB) ([]string, error)
	HasTable(ctx context.Context, db *gorm.DB, table string) (bool, error)
	GetTableColumns(ctx context.Context, db *gorm.DB, table string) ([]TableColumn, error)
}

type Helper struct {
	lock  *sync.RWMutex
	group *singleflight.Group

	dbs map[string]*gorm.DB
}

// NewHelper new a logger helper.
func NewHelper() *Helper {
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	group := &singleflight.Group{}

	return &Helper{
		dbs:   dbs,
		lock:  lock,
		group: group,
	}
}

// GetDB get db connection by key
//
// db := helper.GetDB(ctx)
// db.Table("users").Find(&users)
func (h *Helper) GetDB(options ...DBOption) (*gorm.DB, error) {
	// 修改配置
	config := Apply(options...)

	// 判断是否已经实例化了gorm.DB
	h.lock.RLock()
	if db, ok := h.dbs[config.Dsn]; ok {
		h.lock.RUnlock()
		return db, nil
	}
	h.lock.RUnlock()

	// 没有实例化gorm.DB，那么就要进行实例化操作
	h.lock.Lock()
	defer h.lock.Unlock()

	// use singleflight to avoid multiple goroutines to create the same connection
	do, err, _ := h.group.Do(config.Dsn, func() (interface{}, error) {
		ns := schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			TablePrefix:   fmt.Sprintf("%s_", strings.ReplaceAll(config.Dsn, "-", "")),
		}
		config.NamingStrategy = ns

		// 实例化gorm.DB
		var db *gorm.DB
		var err error
		switch config.DriverType {
		case MySQL:
			db, err = gorm.Open(mysql.Open(config.Dsn), config)
		case PostgreSQL:
			db, err = gorm.Open(postgres.Open(config.Dsn), config)
		case SQLite:
			db, err = gorm.Open(sqlite.Open(config.Dsn), config)
		case SQLServer:
			db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
		case ClickHouse:
			db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
		}
		if err != nil {
			return nil, errors.Wrap(err, "database driver failed")
		}

		// 设置对应的连接池配置
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		// 数据库调优
		if config.ConnMaxIdle > 0 {
			sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
		}
		if config.ConnMaxOpen > 0 {
			sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
		}

		if config.ConnMaxLifetime != "" {
			liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
			if err != nil {
				log.Error(context.Background(), "conn max lift time error", map[string]interface{}{
					"err": err,
				})
			} else {
				sqlDB.SetConnMaxLifetime(liftTime)
			}
		}

		if config.ConnMaxIdleTime != "" {
			idleTime, err := time.ParseDuration(config.ConnMaxIdleTime)
			if err != nil {
				log.Error(context.Background(), "conn max idle time error", map[string]interface{}{
					"err": err,
				})
			} else {
				sqlDB.SetConnMaxIdleTime(idleTime)
			}
		}

		// 挂载到map中，结束配置
		if err == nil {
			h.dbs[config.Dsn] = db
		}

		return db, nil
	})

	if err != nil {
		return nil, err
	}
	return do.(*gorm.DB), nil
}

func (h *Helper) CanConnect(ctx context.Context, db *gorm.DB) (bool, error) {
	sqlDb, err := db.DB()
	if err != nil {
		return false, errors.Wrap(err, "CanConnect")
	}
	if err := sqlDb.Ping(); err != nil {
		return false, errors.Wrap(err, "CanConnect Ping error")
	}
	return true, nil
}

func (h *Helper) GetTables(ctx context.Context, db *gorm.DB) ([]string, error) {
	return db.Migrator().GetTables()
}

func (h *Helper) HasTable(ctx context.Context, db *gorm.DB, table string) (bool, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return false, err
	}
	for _, t := range tables {
		if t == table {
			return true, nil
		}
	}
	return false, nil
}

func (h *Helper) GetTableColumns(ctx context.Context, db *gorm.DB, table string) ([]TableColumn, error) {
	// 执行原始的SQL语句
	var columns []TableColumn
	result := db.Raw("SHOW COLUMNS FROM " + table).Scan(&columns)
	if result.Error != nil {
		// 处理错误
		return nil, result.Error
	}
	return columns, nil
}

// IsDuplicateEntryErr 重复条目, err := db.WithContext(ctx).Create(&user).Error; IsDuplicateEntryErr(err)
func IsDuplicateEntryErr(err error) bool {
	var mysqlErr *mysql2.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == DuplicateEntryErrCode
}

func (h *Helper) CloseAllDB() {
	for name, db := range h.dbs {
		sqlDB, err := db.DB()
		if err != nil {
			log.Errorf("get db instance error: ", err.Error())
			continue
		}

		err = sqlDB.Close()
		if err != nil {
			log.Errorf("close current db error: ", err)
			continue
		}

		// 销毁连接句柄标识
		delete(h.dbs, name)
	}
}

// CloseDbByName 关闭指定name的db engine
func (h *Helper) CloseDbByName(name string) error {
	if db, ok := h.dbs[name]; ok {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		err = sqlDB.Close()
		if err != nil {
			return err
		}

		// 销毁连接句柄标识
		delete(h.dbs, name)
	}

	return errors.New("current db engine not exist")
}
