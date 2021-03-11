package article

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/spaolacci/murmur3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/ioutil"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/common"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/log"
	"mcm-api/pkg/media"
	"mcm-api/pkg/queue"
	"time"
)

type Service struct {
	cfg          *config.Config
	repository   *repository
	mediaService media.Service
	queue        queue.Queue
}

func InitializeService(
	cfg *config.Config,
	repository *repository,
	mediaService media.Service,
	queue queue.Queue,
) *Service {
	return &Service{
		cfg:          cfg,
		repository:   repository,
		mediaService: mediaService,
		queue:        queue,
	}
}

func (s Service) FindById(ctx context.Context, id int) (*ArticleRes, error) {
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "article not found", err)
		}
		return nil, err
	}
	return s.mapArticleToRes(entity), nil
}

func (s Service) Create(ctx context.Context, req *ArticleReq) (*ArticleRes, error) {
	file, err := s.mediaService.GetFile(ctx, req.Link)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fileHash := hash(fileContent)
	entity, err := s.repository.Create(ctx, &Entity{
		Title:       req.Title,
		Description: req.Description,
		Versions: []*Version{
			{
				Hash:         fileHash,
				LinkOriginal: req.Link,
				LinkPdf:      "",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	user, err := enforcer.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	s.addToQueue(user, entity.Versions[0])
	return s.mapArticleToRes(entity), nil
}

func (s Service) addToQueue(user *enforcer.LoggedInUser, entity *Version) {
	go func() {
		ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
		defer cancelFunc()
		er := s.queue.Add(ctxTimeout, &queue.Message{
			Topic: queue.ArticleUploaded,
			Data: &queue.ArticleUploadedPayload{
				ArticleId: entity.Id,
				Link:      entity.LinkOriginal,
				User:      *user,
			},
		})
		if er != nil {
			log.Logger.Error("add message to queue failed", zap.Error(er))
		}
	}()
}

func (s Service) Update(ctx context.Context, articleId int, req ArticleReq) (*ArticleRes, error) {
	entity, err := s.repository.FindById(ctx, articleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "article not found", err)
		}
		return nil, err
	}
	entity.Title = req.Title
	entity.Description = req.Description
	entity, err = s.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}
	if req.Link != "" {
		_, err = s.CreateVersion(ctx, articleId, req.Link)
		if err != nil {
			return nil, err
		}
	}
	return s.mapArticleToRes(entity), nil
}

func (s Service) CreateVersion(ctx context.Context, articleId int, link string) (*Version, error) {
	file, err := s.mediaService.GetFile(ctx, link)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fileHash := hash(fileContent)
	latestVersionOfArticle, err := s.repository.GetLatestVersionOfArticle(ctx, articleId)
	if err != nil {
		return nil, err
	}
	if fileHash == latestVersionOfArticle.Hash {
		return nil, apperror.New(apperror.ErrConflict, "duplicate article version", nil)
	}
	version, err := s.repository.CreateVersion(ctx, &Version{
		Hash:         fileHash,
		ArticleId:    articleId,
		LinkOriginal: link,
	})
	if err != nil {
		return nil, err
	}
	user, err := enforcer.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	s.addToQueue(user, version)
	return version, nil
}

func (s Service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func hash(input []byte) string {
	new128 := murmur3.New128()
	_, _ = new128.Write(input)
	return hex.EncodeToString(new128.Sum(nil))
}

func (s Service) mapArticleToRes(a *Entity) *ArticleRes {
	return &ArticleRes{
		Id:          a.Id,
		Title:       a.Title,
		Description: a.Description,
		Versions:    s.mapVersionsToRes(a.Versions...),
		TrackTime: common.TrackTime{
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
	}
}

func (s Service) mapVersionsToRes(vs ...*Version) []*VersionRes {
	var results []*VersionRes
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	for _, v := range vs {
		res := &VersionRes{
			Id:           v.Id,
			Hash:         v.Hash,
			ArticleId:    v.ArticleId,
			LinkOriginal: v.LinkOriginal,
			LinkPdf:      v.LinkPdf,
			CreatedAt:    v.CreatedAt,
		}

		url, err := s.mediaService.GetUrl(ctx, v.LinkOriginal)
		if err != nil {
			log.Logger.Error("error get url", zap.Error(err))
		}
		res.LinkOriginalCdn = url
		if v.LinkPdf != "" {
			url, err = s.mediaService.GetUrl(ctx, v.LinkPdf)
			if err != nil {
				log.Logger.Error("error get url", zap.Error(err))
			}
			res.LinkPdfCdn = url
		}
		results = append(results, res)
	}
	return results
}

func (s Service) UpdateLinkPdfForVersion(ctx context.Context, id int, key string) error {
	entity, err := s.repository.FindVersionById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.ErrNotFound, "article version not found", err)
		}
		return err
	}
	entity.LinkPdf = key
	return s.repository.UpdateVersion(ctx, entity)
}
