package generics

import (
	"context"

	"gorm.io/gorm"
)

type Model[E any] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}

type Repository[M Model[E], E any] struct {
	db *gorm.DB
}

func NewRepository[M Model[E], E any](db *gorm.DB) *Repository[M, E] {
	return &Repository[M, E]{
		db: db,
	}
}

func (r *Repository[M, E]) Insert(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

func (r *Repository[M, E]) Delete(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)
	err := r.db.WithContext(ctx).Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository[M, E]) DeleteByID(ctx context.Context, id any) error {
	var start M
	err := r.db.WithContext(ctx).Delete(&start, &id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[M, E]) Update(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

func (r *Repository[M, E]) FindByID(ctx context.Context, id any) (E, error) {
	var model M
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return *new(E), err
	}

	return model.ToEntity(), nil
}

func (r *Repository[M, E]) Find(ctx context.Context, specifications ...Specification) ([]E, error) {
	return r.FindWithLimit(ctx, -1, -1, specifications...)
}

func (r *Repository[M, E]) Count(ctx context.Context, specifications ...Specification) (i int64, err error) {
	model := new(M)
	err = r.getPreWarmDbForSelect(ctx, specifications...).Model(model).Count(&i).Error
	return
}

func (r *Repository[M, E]) getPreWarmDbForSelect(ctx context.Context, specification ...Specification) *gorm.DB {
	dbPrewarm := r.db.WithContext(ctx)
	for _, s := range specification {
		dbPrewarm = dbPrewarm.Where(s.GetQuery(), s.GetValues()...)
	}
	return dbPrewarm
}
func (r *Repository[M, E]) FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]E, error) {
	var models []M

	dbPrewarm := r.getPreWarmDbForSelect(ctx, specifications...)
	err := dbPrewarm.Limit(limit).Offset(offset).Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make([]E, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToEntity())
	}

	return result, nil
}

func (r *Repository[M, E]) FindAll(ctx context.Context) ([]E, error) {
	return r.FindWithLimit(ctx, -1, -1)
}
