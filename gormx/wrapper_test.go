package gormx

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
	"testing"
)

type Database struct {
	*Helper
}

type myUser struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(200)"`
}

func (myUser) TableName() string {
	return "user"
}

func TestDatabase_OpenDB(t *testing.T) {
	ctx := context.Background()
	h := &Database{NewHelper()}

	db, err := h.GetDB(WithConfig(func(options *DBConfig) {
		options.DriverType = MySQL
		options.Dsn = "root:123456@tcp(localhost:13306)/webserver_db?charset=utf8mb4&parseTime=True&loc=Local"
		options.DisableForeignKeyConstraintWhenMigrating = true

		tag := options.DriverType.String()
		if tag == "mysql" {
			t.Log("mysql")
		}
	}))

	if err != nil {
		t.Fatal(err)
	}

	// 检测数据库是否可以连接
	canConnected, err := h.CanConnect(ctx, db)
	if err != nil || !canConnected {
		fmt.Println("数据库连接失败，请检查配置")
		t.Fatal(err)
	}

	tables, err := h.GetTables(ctx, db)
	if err != nil {
		fmt.Println("获取数据库表格失败")
		t.Fatal(err)
	}
	t.Log(tables)

	table := ""

	hasTable, err := h.HasTable(ctx, db, table)
	if err != nil {
		t.Errorf("数据库连接失败，表格 %v, 错误 %v", table, err)
	}

	if hasTable == false {
		t.Errorf("表格 %v 不存在", table)
	}

	// 获取所有字段
	columns, err := h.GetTableColumns(ctx, db, table)
	if err != nil {
		t.Errorf("获取表格 %v 列表字段失败: %v", table, err)
	}

	t.Log(columns)

	tableLower := strings.ToLower(table)

	t.Log(tableLower)
}

func TestShortConnect(t *testing.T) {
	getDb := func() (*gorm.DB, error) {
		ctx := context.Background()
		h := &Database{NewHelper()}

		db, err := h.GetDB(WithConfig(func(options *DBConfig) {
			options.DriverType = MySQL
			options.Dsn = "root:123456@tcp(localhost:13306)/webserver_db?charset=utf8mb4&parseTime=True&loc=Local"
			options.DisableForeignKeyConstraintWhenMigrating = true

			tag := options.DriverType.String()
			if tag == "mysql" {
				t.Log("mysql")
			}
		}))

		if err != nil {
			t.Fatal(err)
		}

		// 检测数据库是否可以连接
		canConnected, err := h.CanConnect(ctx, db)
		if err != nil || !canConnected {
			fmt.Println("数据库连接失败，请检查配置")
			t.Fatal(err)
		}

		return db, nil
	}

	// 这里我设置了db max_connections最大连接为1000
	var wg sync.WaitGroup
	// var maxConnections = 30
	var maxConnections = 1000
	// var maxConnections = 2000
	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			db, err := getDb()
			if err != nil {
				log.Println("get db error: ", err.Error())
				return
			}

			// 关闭数据句柄
			defer func() {
				sqlDB, err := db.DB()
				if err != nil {
					log.Println("get db instance error: ", err.Error())
					return
				}

				_ = sqlDB.Close()
			}()

			user := &myUser{}
			db.Where("name = ?", "hello").First(user)
			log.Println(user)
		}()
	}

	wg.Wait()
	log.Println("test success")
}
