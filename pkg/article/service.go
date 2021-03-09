package article

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/spaolacci/murmur3"
	"gorm.io/gorm"
	"io/ioutil"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/common"
	"mcm-api/pkg/media"
)

type Service struct {
	cfg          *config.Config
	repository   *repository
	mediaService media.Service
}

func InitializeService(
	cfg *config.Config,
	repository *repository,
	mediaService media.Service,
) *Service {
	return &Service{
		cfg:          cfg,
		repository:   repository,
		mediaService: mediaService,
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
	return mapArticleToRes(entity), nil
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
		Versions: []Version{
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
	return mapArticleToRes(entity), nil
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
	return mapArticleToRes(entity), nil
}

func (s Service) CreateVersion(ctx context.Context, articleId int, link string) (*Version, error) {
	file, err := s.mediaService.GetFile(ctx, link)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return nil, err
	}
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
	return version, nil
}

func hash(input []byte) string {
	new128 := murmur3.New128()
	_, _ = new128.Write(input)
	return hex.EncodeToString(new128.Sum(nil))
}

func mapArticleToRes(a *Entity) *ArticleRes {
	return &ArticleRes{
		Id:          a.Id,
		Title:       a.Title,
		Description: a.Description,
		Versions:    a.Versions,
		TrackTime: common.TrackTime{
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
	}
}
