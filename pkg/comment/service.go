package comment

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/common"
	"mcm-api/pkg/user"
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
	entities, nextCursor, err := s.repository.FindCursor(ctx, query)
	if err != nil {
		return nil, err
	}

	var res []*CommentRes
	for _, v := range entities {
		res = append(res, mapEntityToRes(v))
	}

	return common.NewCursorResponse(res, nextCursor), nil
}

func (s Service) FindById(ctx context.Context, id string) (*CommentRes, error) {
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "comment not found", err)
		}
		return nil, err
	}
	return mapEntityToRes(entity), nil
}

func (s Service) Create(ctx context.Context, body *CommentCreateReq) (*CommentRes, error) {
	u, err := common.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	err = body.Validate()
	if err != nil {
		return nil, err
	}
	entity, err := s.repository.Create(ctx, &Entity{
		UserId:         u.Id,
		ContributionId: body.ContributionId,
		Content:        body.Content,
	})
	if err != nil {
		return nil, err
	}
	return mapEntityToRes(entity), nil
}

func (s Service) Update(ctx context.Context, id string, body *CommentUpdateReq) (*CommentRes, error) {
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "comment not found", err)
		}
		return nil, err
	}
	entity.Content = body.Content
	entity, err = s.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapEntityToRes(entity), nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func mapEntityToRes(entity *Entity) *CommentRes {
	return &CommentRes{
		Id: entity.Id,
		User: user.UserResponse{
			Id:        entity.User.Id,
			Name:      entity.User.Name,
			Email:     entity.User.Email,
			FacultyId: entity.User.FacultyId,
			Role:      entity.User.Role,
			TrackTime: common.TrackTime{
				CreatedAt: entity.User.CreatedAt,
				UpdatedAt: entity.User.UpdatedAt,
			},
		},
		Content: entity.Content,
		Edited:  !entity.CreatedAt.Equal(entity.UpdatedAt),
		TrackTime: common.TrackTime{
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
	}
}
