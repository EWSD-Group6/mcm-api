package comment

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/common"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/user"
)

type Service struct {
	cfg                 *config.Config
	repository          *repository
	contributionService *contribution.Service
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
	u, _ := enforcer.GetLoggedInUser(ctx)
	contrib, err := s.contributionService.FindById(ctx, query.ContributionId)
	if err != nil {
		return nil, err
	}
	err = canCommentOnContribution(u, contrib)
	if err != nil {
		return nil, err
	}
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
	u, err := enforcer.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	err = body.Validate()
	if err != nil {
		return nil, err
	}
	ctb, err := s.contributionService.FindById(ctx, body.ContributionId)
	if err != nil {
		return nil, err
	}
	err = canCommentOnContribution(u, ctb)
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

func canCommentOnContribution(user *enforcer.LoggedInUser, contrib *contribution.ContributionRes) error {
	if user.Role == enforcer.Student &&
		contrib.User.Id != user.Id {
		return apperror.New(apperror.ErrForbidden,
			"you are not owner of this contribution", nil)
	}
	if user.Role == enforcer.MarketingCoordinator &&
		contrib.User.FacultyId != user.FacultyId {
		return apperror.New(
			apperror.ErrForbidden,
			"you are not in same faculty with contribution", nil)
	}
	return nil
}

func (s Service) Update(ctx context.Context, id string, body *CommentUpdateReq) (*CommentRes, error) {
	u, _ := enforcer.GetLoggedInUser(ctx)
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "comment not found", err)
		}
		return nil, err
	}
	if entity.UserId != u.Id {
		return nil, apperror.New(apperror.ErrForbidden, "not your comment", nil)
	}
	entity.Content = body.Content
	entity, err = s.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapEntityToRes(entity), nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	u, _ := enforcer.GetLoggedInUser(ctx)
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.ErrNotFound, "comment not found", err)
		}
		return err
	}
	if entity.UserId != u.Id {
		return apperror.New(apperror.ErrForbidden, "not your comment", nil)
	}
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
