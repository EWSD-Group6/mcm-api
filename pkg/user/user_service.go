package user

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/common"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/faculty"
	"mcm-api/pkg/log"
)

type Service struct {
	cfg            *config.Config
	repository     *repository
	facultyService *faculty.Service
}

func InitializeService(
	cfg *config.Config,
	repository *repository,
	facultyService *faculty.Service,
) *Service {
	return &Service{
		cfg:            cfg,
		repository:     repository,
		facultyService: facultyService,
	}
}

func (s *Service) FindById(ctx context.Context, id int) (*UserResponse, error) {
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "user not found", err)
		}
		return nil, err
	}
	return mapEntityToResponse(entity), nil
}

func (s *Service) FindByEmailAndPassword(ctx context.Context, email string, password string) (*UserResponse, error) {
	entity, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "user not found", err)
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(password))
	if err != nil {
		return nil, apperror.New(apperror.ErrInvalid, "wrong password", err)
	}
	return mapEntityToResponse(entity), nil
}

func (s *Service) CreateUser(ctx context.Context, req *UserCreateReq) (*UserResponse, error) {
	// common validate
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	// validate duplicate
	entity, err := s.repository.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, apperror.New(apperror.ErrConflict, "duplicate email", err)
	}
	entity = &Entity{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	// validate and set faculty
	if req.Role == enforcer.MarketingCoordinator ||
		req.Role == enforcer.Student ||
		req.Role == enforcer.Guest {
		_, err = s.facultyService.FindById(ctx, *req.FacultyId)
		if err != nil {
			return nil, apperror.New(apperror.ErrInvalid, "invalid faculty", err)
		}
		entity.FacultyId = req.FacultyId
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}
	entity.Password = string(hashedPassword)

	// save
	err = s.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapEntityToResponse(entity), nil
}

func (s *Service) CreateDefaultAdmin(ctx context.Context) error {
	_, err := s.repository.FindByEmail(ctx, s.cfg.AdminEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Logger.Info("skip create default admin")
		return nil
	}
	log.Logger.Info("init admin account", zap.String("email", s.cfg.AdminEmail))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.cfg.AdminPassword), 10)
	if err != nil {
		return err
	}
	return s.repository.Create(ctx, &Entity{
		Name:     "Administrators",
		Email:    s.cfg.AdminEmail,
		Password: string(hashedPassword),
		Role:     enforcer.Administrator,
	})
}

func (s *Service) Find(ctx context.Context, user *enforcer.LoggedInUser, query *UserIndexQuery) (*common.PaginateResponse, error) {
	entities, count, err := s.repository.FindAndCount(ctx, query)
	if err != nil {
		return nil, err
	}
	dtos := mapEntitiesToResponse(entities)
	return common.NewPaginateResponse(dtos, count, query.Page, query.GetLimit()), nil
}

func (s *Service) GetAllUserOfFaculty(ctx context.Context, role enforcer.Role, facultyId int) ([]*Entity, error) {
	return s.repository.FindAllUserOfFaculty(ctx, role, facultyId)
}

func (s *Service) Update(ctx context.Context, id int, req *UserUpdateReq) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	entity, err := s.repository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "user not found", err)
		}
		return nil, err
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), 10)
		if err != nil {
			return nil, err
		}
		entity.Password = string(hashedPassword)
	}
	if req.Role != nil {
		entity.Role = *req.Role
	}
	if req.Email != nil {
		entity.Email = *req.Email
	}
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.FacultyId != nil {
		entity.FacultyId = req.FacultyId
	}

	entity, err = s.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapEntityToResponse(entity), nil
}

func mapEntitiesToResponse(entity []*Entity) []*UserResponse {
	var result []*UserResponse
	for i := range entity {
		result = append(result, mapEntityToResponse(entity[i]))
	}
	return result
}

func mapEntityToResponse(entity *Entity) *UserResponse {
	return &UserResponse{
		Id:        entity.Id,
		Name:      entity.Name,
		Role:      entity.Role,
		Email:     entity.Email,
		FacultyId: entity.FacultyId,
		TrackTime: common.TrackTime{
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
	}
}
