package comment

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

func (s Service) Find(ctx context.Context, query *IndexQuery) (*common.CursorResponse, error) {
	s.repository.Find(ctx, query)
	return nil, nil
}

func (s Service) FindById(ctx context.Context, id int) (*CommentRes, error) {
	return nil, nil
}

func (s Service) Create(ctx context.Context, body *CommentCreateReq) (*CommentRes, error) {
	return nil, nil
}

func (s Service) Update(ctx context.Context, id int, body *CommentUpdateReq) (*CommentRes, error) {
	return nil, nil
}

func (s Service) Delete(ctx context.Context, id int) error {
	return nil
}
