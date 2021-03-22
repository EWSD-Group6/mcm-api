package systemdata

import (
	"context"
	"mcm-api/config"
	"mcm-api/pkg/common"
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

func (s Service) Find(ctx context.Context) (*common.PaginateResponse, error) {
	return nil, nil
}

func (s Service) Update(ctx context.Context, id string, body *DataUpdateReq) error {
	return nil
}
