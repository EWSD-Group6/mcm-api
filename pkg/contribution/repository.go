package contribution

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

func (r repository) FindById(ctx context.Context, id int) (*Entity, error) {
	result := new(Entity)
	db := r.db.WithContext(ctx).Preload("User").First(result, id)
	return result, db.Error
}

func (r repository) Find(ctx context.Context, query *IndexQuery) ([]*Entity, error) {
	var results []*Entity
	r.db.WithContext(ctx)
	r.db.Limit(query.GetLimit())
	r.db.Offset(query.GetOffSet())
	db := r.db.Find(results)
	return results, db.Error
}

func (r repository) Create(ctx context.Context, entity *Entity) (*Entity, error) {
	db := r.db.WithContext(ctx).Create(entity)
	return entity, db.Error
}

func (r repository) Update(ctx context.Context, entity *Entity) (*Entity, error) {
	db := r.db.WithContext(ctx).Save(entity)
	return entity, db.Error
}

func (r repository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Entity{}, id).Error
}

func (r repository) FindAndCount(ctx context.Context, query *IndexQuery) ([]*Entity, int64, error) {
	var entities []*Entity
	builder := r.db.WithContext(ctx).Model(&Entity{})
	if query.Status != "" {
		builder.Where("status = ?", query.Status)
	}
	if query.FacultyId != nil {
		builder.Joins("join users on users.id = contributions.user_id").
			Where("users.faculty_id = ?", query.FacultyId)
	}
	if query.StudentId != nil {
		builder.Where("user_id = ?", query.StudentId)
	}
	if query.ContributionSessionId != nil {
		builder.Where("contribute_session_id = ?", query.ContributionSessionId)
	}
	var count int64
	result := builder.Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	builder.Offset(query.GetOffSet()).Limit(query.GetLimit())
	result = builder.Preload("User").Preload("Article").Find(&entities)
	return entities, count, result.Error
}

func (r repository) GetImagesById(ctx context.Context, id int) ([]*ImageEntity, error) {
	var entities []*ImageEntity
	result := r.db.WithContext(ctx).Where("contribution_id = ?", id).Find(&entities)
	return entities, result.Error
}

func (r repository) GetAllAcceptedContributions(ctx context.Context, contributeSessionId int) ([]*Entity, error) {
	var entities []*Entity
	result := r.db.WithContext(ctx).
		Where("status = ? and contribute_session_id = ?", Accepted, contributeSessionId).
		Find(&entities)
	return entities, result.Error
}
