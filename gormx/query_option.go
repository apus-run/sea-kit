package gormx

import (
	"fmt"

	"gorm.io/gorm"
)

type QueryOption func(db *gorm.DB) *gorm.DB

func WithID(id uint) QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("id = ?", id)
	}
}

func WithIDs(ids []uint) QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("id IN ?", ids)
	}
}

func WithStatus(status int32) QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("status = ?", status)
	}
}

func WithPageLimit(offset, limit int) QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Offset(offset).Limit(limit)
	}
}

func WithAsc() QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Order("created_at ASC")
	}
}

func WithDesc() QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Order("created_at DESC")
	}
}

func WithFuzzyName(name string) QueryOption {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where(fmt.Sprintf("name like %q", ("%" + name + "%")))
	}
}

// 使用
// type userDAO interface {
// 	GetUsers(ctx context.Context, opts ...QueryOption) ([]*model.User, error)
// }

// type UserDAO struct {
// 	client *data.DB
// }
// func (u *UserDAO) GetUsers(ctx context.Context, opts ...QueryOption) ([]*model.User, error) {
// 	db := t.client.DB.WithContext(ctx)
// 	for _, opt := range opts {
// 		db = opt(db)
// 	}

// 	var users []*model.User
// 	return users, db.Model(&model.User{}).Scan(&users).Error
// }

// tasks, err := GetUsers(ctx, dao.WithAsc(), dao.WithStatus(1))
