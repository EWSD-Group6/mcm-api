package systemdata

import (
	"context"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func InitializeRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r repository) Find(ctx context.Context) ([]*Entity, error) {
	var results []*Entity
	r.db.WithContext(ctx)
	db := r.db.Find(&results)
	return results, db.Error
}

func (r repository) Update(ctx context.Context, entity *Entity) (*Entity, error) {
	db := r.db.WithContext(ctx).Save(entity)
	return entity, db.Error
}

func (r repository) FindById(ctx context.Context, key string) (*Entity, error) {
	entity := &Entity{}
	db := r.db.WithContext(ctx).First(entity, "key = ?", key)
	if db.Error != nil {
		return nil, db.Error
	}
	return entity, nil
}
