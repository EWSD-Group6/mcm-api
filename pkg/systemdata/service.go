package systemdata

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
)

type Service struct {
	cfg        *config.Config
	repository *repository
}

func InitializeService(
	cfg *config.Config,
	repository *repository,
) *Service {
	return &Service{
		cfg:        cfg,
		repository: repository,
	}
}

func (s Service) Find(ctx context.Context) ([]*DataRes, error) {
	entities, err := s.repository.Find(ctx)
	if err != nil {
		return nil, err
	}
	return mapEntitiesToRes(entities), nil
}

func (s Service) Update(ctx context.Context, key string, body *DataUpdateReq) error {
	entity, err := s.repository.FindById(ctx, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.ErrNotFound, "system data not found", err)
		}
		return err
	}
	entity.Value = body.Value
	_, err = s.repository.Update(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func mapEntitiesToRes(entities []*Entity) []*DataRes {
	var res []*DataRes
	for _, v := range entities {
		res = append(res, &DataRes{
			Key:       v.Key,
			Value:     v.Value,
			Type:      v.Type,
			UpdatedAt: v.UpdatedAt,
		})
	}
	return res
}
