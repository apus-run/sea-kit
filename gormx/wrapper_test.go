package gormx

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type Database struct {
	*Helper
}

func TestDatabase_OpenDB(t *testing.T) {
	ctx := context.Background()
	h := &Database{NewHelper()}

	db, err := h.GetDB(WithGormConfig(func(options *DBConfig) {
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
