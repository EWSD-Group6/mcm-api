package article

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
	db := r.db.WithContext(ctx).First(result, id)
	return result, db.Error
}

func (r repository) Create(ctx context.Context, entity *Entity) (*Entity, error) {
	result := r.db.WithContext(ctx).Create(entity)
	return entity, result.Error
}

func (r repository) CreateVersion(ctx context.Context, entity *Version) (*Version, error) {
	result := r.db.WithContext(ctx).Create(entity)
	return entity, result.Error
}

func (r repository) Update(ctx context.Context, entity *Entity) (*Entity, error) {
	result := r.db.WithContext(ctx).Save(entity)
	return entity, result.Error
}

func (r repository) GetLatestVersionOfArticle(ctx context.Context, articleId int) (*Version, error) {
	entity := new(Version)
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Where("article_id = ?", articleId).
		Find(entity)
	return entity, result.Error
}
