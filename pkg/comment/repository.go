package comment

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

func (r repository) FindById(ctx context.Context, id string) (*Entity, error) {
	result := new(Entity)
	db := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(result)
	return result, db.Error
}

func (r repository) FindCursor(ctx context.Context, query *IndexQuery) ([]*Entity, *CursorPayload, error) {
	builder := r.db.Debug().WithContext(ctx).
		Preload("User").
		Limit(query.GetLimit()).
		Order("created_at ASC").
		Where("contribution_id = ?", query.ContributionId)
	if query.Next != "" {
		var nextPayload CursorPayload
		err := query.GetNext(&nextPayload)
		if err != nil {
			return nil, nil, err
		}
		builder.Where("created_at >= ? and id != ?", nextPayload.CreatedAt, nextPayload.Id)
	}
	var entities []*Entity
	result := builder.Find(&entities)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	var nextCursor *CursorPayload
	if len(entities) > 0 {
		lastEntity := entities[len(entities)-1]
		nextCursor = &CursorPayload{
			Id:        lastEntity.Id,
			CreatedAt: lastEntity.CreatedAt,
		}
	}
	return entities, nextCursor, nil
}

func (r repository) Create(ctx context.Context, entity *Entity) (*Entity, error) {
	db := r.db.WithContext(ctx).Create(entity)
	return entity, db.Error
}

func (r repository) Update(ctx context.Context, entity *Entity) (*Entity, error) {
	db := r.db.WithContext(ctx).Save(entity)
	return entity, db.Error
}

func (r repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Entity{}).Error
}
