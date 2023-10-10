package gormx

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/apus-run/sea-kit/pagination"
)

func LikeByBigIntColumn(db *gorm.DB, column, val string) *gorm.DB {
	if len(column) > 0 {
		return db.Where(fmt.Sprintf("CAST(%s as char) like ?", column), "%"+val+"%")
	}
	return db
}

func OptionWhereString(db *gorm.DB, whereStr string, fuzz string) *gorm.DB {
	if len(fuzz) > 0 {
		if strings.Contains(whereStr, " like ") {
			return db.Where(whereStr, "%"+fuzz+"%")
		} else {
			return db.Where(whereStr, fuzz)
		}
	}
	return db
}
func OptionWhereBool(db *gorm.DB, whereStr string, b *bool) *gorm.DB {
	if b != nil {
		return db.Where(whereStr, b)
	}
	return db
}
func OptionWhereId(db *gorm.DB, whereStr string, id int64) *gorm.DB {
	if id > 0 {
		return db.Where(whereStr, id)
	}
	return db
}
func OptionWhereInt(db *gorm.DB, whereStr string, val int64) *gorm.DB {
	if val != 0 {
		return db.Where(whereStr, val)
	}
	return db
}
func OptionWhereInts(db *gorm.DB, whereStr string, val []int64) *gorm.DB {
	if len(val) > 0 {
		return db.Where(whereStr, val)
	}
	return db
}

func CreateTimeRange(startTime, endTime *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if startTime != nil {
			db = db.Where("created_at >= ?", *startTime)
		}
		if endTime != nil {
			db = db.Where("created_at <= ?", *endTime)
		}
		return db
	}
}

func Paginate(page *pagination.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == nil {
			return db
		}
		return db.Offset(page.Offset()).Limit(page.Limit())
	}
}

// AvgRaw 用于处理 avg() 为 null 的情况
func AvgRaw(column, alias string) string {
	return fmt.Sprintf("IFNULL(avg(%s),0) as %s", column, alias)
}

// CaseWhenNull time range
func CaseWhenNull(column, alias string, defaultValue interface{}) string {
	var fmtStr = "case WHEN avg( %s )  IS NULL THEN %v ELSE avg( %s ) END %s"
	if _, ok := defaultValue.(int); ok {
		fmtStr = "case WHEN avg( %s )  IS NULL THEN %d ELSE avg( %s ) END %s"
	} else if _, ok := defaultValue.(float64); ok {
		fmtStr = "case WHEN avg( %s )  IS NULL THEN %.2f ELSE avg( %s ) END %s"
	}

	return fmt.Sprintf(fmtStr, column, defaultValue, column, alias)
}
