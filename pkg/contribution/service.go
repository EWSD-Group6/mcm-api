package contribution

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/article"
	"mcm-api/pkg/common"
	"mcm-api/pkg/contributesession"
	"mcm-api/pkg/media"
	"mcm-api/pkg/queue"
)

type Service struct {
	cfg                      *config.Config
	repository               *repository
	queue                    queue.Queue
	contributeSessionService *contributesession.Service
	articleService           *article.Service
	mediaService             media.Service
}

func InitializeService(
	cfg *config.Config,
	repository *repository,
	queue queue.Queue,
	cs *contributesession.Service,
	articleService *article.Service,
	mediaService media.Service,
) *Service {
	return &Service{
		queue:                    queue,
		cfg:                      cfg,
		repository:               repository,
		contributeSessionService: cs,
		articleService:           articleService,
		mediaService:             mediaService,
	}
}

func (s Service) Find(ctx context.Context, query *IndexQuery) (*common.PaginateResponse, error) {
	loggedInUser, err := common.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	var result []*Entity
	var count int64
	switch loggedInUser.Role {
	case common.Administrator:
		result, count, err = s.repository.FindAndCount(ctx, &IndexQuery{
			PaginateQuery:         query.PaginateQuery,
			FacultyId:             query.FacultyId,
			StudentId:             query.StudentId,
			ContributionSessionId: query.ContributionSessionId,
			Status:                "",
		})
		if err != nil {
			return nil, err
		}
		break
	case common.MarketingCoordinator:
		result, count, err = s.repository.FindAndCount(ctx, &IndexQuery{
			PaginateQuery:         query.PaginateQuery,
			FacultyId:             loggedInUser.FacultyId,
			StudentId:             query.StudentId,
			ContributionSessionId: query.ContributionSessionId,
			Status:                query.Status,
		})
		if err != nil {
			return nil, err
		}
		break
	case common.Student:
		result, count, err = s.repository.FindAndCount(ctx, &IndexQuery{
			PaginateQuery:         query.PaginateQuery,
			FacultyId:             loggedInUser.FacultyId,
			StudentId:             &loggedInUser.Id,
			ContributionSessionId: query.ContributionSessionId,
			Status:                query.Status,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, apperror.New(apperror.ErrForbidden, "", nil)
	}

	return common.NewPaginateResponse(
		mapManyContributionToRes(result),
		count,
		query.Page,
		query.GetLimit(),
	), nil
}

func (s Service) FindById(ctx context.Context, id int) (*ContributionRes, error) {
	entity, err := s.findById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapContributionToRes(entity), nil
}

func (s Service) findById(ctx context.Context, id int) (*Entity, error) {
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "contribution not found", err)
		}
		return nil, err
	}
	return entity, nil
}

func (s Service) Create(ctx context.Context, body *ContributionCreateReq) (*ContributionRes, error) {
	loggedInUser, err := common.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	session, err := s.contributeSessionService.GetCurrentSession(ctx)
	if err != nil {
		return nil, err
	}
	var a *article.ArticleRes
	if body.Article != nil {
		a, err = s.articleService.Create(ctx, &article.ArticleReq{
			Title:       body.Article.Title,
			Description: body.Article.Description,
			Link:        body.Article.Link,
		})
		if err != nil {
			return nil, err
		}
	}
	entity := &Entity{
		UserId:              loggedInUser.Id,
		ContributeSessionId: session.Id,
		Status:              Reviewing,
		Images:              mapImageReqToEntity(body.Images...),
	}
	if a != nil {
		entity.ArticleId = &a.Id
	}
	entity, err = s.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapContributionToRes(entity), nil
}

func (s Service) Update(ctx context.Context, id int, body *ContributionUpdateReq) (*ContributionRes, error) {
	entity, err := s.findById(ctx, id)
	if err != nil {
		return nil, err
	}
	if body.Article != nil {
		_, err = s.articleService.Update(ctx, *entity.ArticleId, article.ArticleReq{
			Title:       body.Article.Title,
			Description: body.Article.Description,
			Link:        body.Article.Link,
		})
		if err != nil {
			return nil, err
		}
	}
	if body.Images != nil {
		entity.Images = mapImageReqToEntity(body.Images...)
		_, err = s.repository.Update(ctx, entity)
		if err != nil {
			return nil, err
		}
	}
	return mapContributionToRes(entity), nil
}

func (s Service) Delete(ctx context.Context, id int) error {
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	err = s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return s.articleService.Delete(ctx, *entity.ArticleId)
}

func (s Service) GetImages(ctx context.Context, id int) ([]*ImageRes, error) {
	entities, err := s.repository.GetImagesById(ctx, id)
	if err != nil {
		return nil, err
	}
	var res []*ImageRes
	for _, v := range entities {
		res = append(res, &ImageRes{
			Key:   v.Key,
			Title: v.Title,
			Link:  s.mediaService.GetImageLink(v.Key),
		})
	}
	return res, nil
}

func mapImageReqToEntity(images ...ImageCreateReq) []ImageEntity {
	var result []ImageEntity
	for i := range images {
		result = append(result, ImageEntity{
			Key:   images[i].Key,
			Title: images[i].Title,
		})
	}
	return result
}

func mapContributionToRes(c *Entity) *ContributionRes {
	return &ContributionRes{
		Id: c.Id,
		User: UserRes{
			Id:        c.User.Id,
			Name:      c.User.Name,
			Email:     c.User.Email,
			FacultyId: c.User.FacultyId,
			Role:      c.User.Role,
		},
		ContributeSessionId: c.ContributeSessionId,
		ArticleId:           c.ArticleId,
		Status:              c.Status,
		TrackTime: common.TrackTime{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
	}
}

func mapManyContributionToRes(entities []*Entity) []*ContributionRes {
	var result []*ContributionRes
	for _, v := range entities {
		result = append(result, mapContributionToRes(v))
	}
	return result
}
