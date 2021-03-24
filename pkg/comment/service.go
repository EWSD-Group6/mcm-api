package comment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/common"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/log"
	"mcm-api/pkg/user"
)

type Service struct {
	cfg                 *config.Config
	repository          *repository
	contributionService *contribution.Service
	redis               *redis.Client
}

func InitializeService(
	cfg *config.Config,
	repository *repository,
	redis *redis.Client,
	contributionService *contribution.Service,
) *Service {
	return &Service{
		cfg:                 cfg,
		repository:          repository,
		redis:               redis,
		contributionService: contributionService,
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
		ContributionId: ctb.Id,
		Content:        body.Content,
	})
	if err != nil {
		return nil, err
	}
	res := mapEntityToRes(entity)
	s.publishComment(ctx, ctb.Id, res)
	return res, nil
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

func (s Service) StreamingComment(ctx context.Context, contributionId int, channel chan CommentRes) error {
	_, err := s.contributionService.FindById(ctx, contributionId)
	if err != nil {
		return err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     s.cfg.RedisAddr,
		Password: s.cfg.RedisPassword, // no password set
		DB:       s.cfg.RedisDb,       // use default DB
	})
	pubSub := rdb.Subscribe(ctx, generateCommentChannelName(contributionId))
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case message := <-pubSub.Channel():
			comment := &CommentRes{}
			err = json.Unmarshal([]byte(message.Payload), comment)
			if err != nil {
				log.Logger.Error("error parsing json from comment channel", zap.Error(err))
			}
			channel <- *comment
		}
	}
	return nil
}

func (s Service) publishComment(ctx context.Context, contributionId int, res *CommentRes) {
	bytes, err := json.Marshal(res)
	if err != nil {
		log.Logger.Error("marshal json failed", zap.Error(err))
	}
	s.redis.Publish(ctx, generateCommentChannelName(contributionId), string(bytes))
}

func generateCommentChannelName(contributionId int) string {
	return fmt.Sprintf("contribution:%v:comment-channel", contributionId)
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
