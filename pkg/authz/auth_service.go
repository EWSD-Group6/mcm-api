package authz

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"mcm-api/config"
	"mcm-api/pkg/apperror"
	"mcm-api/pkg/user"
	"strconv"
	"time"
)

// access token ttl, unit: hours
const accessTokenTtl = 168

type Service struct {
	config      *config.Config
	userService *user.Service
}

func InitializeAuthService(config *config.Config, userService *user.Service) *Service {
	return &Service{
		config:      config,
		userService: userService,
	}
}

func (s Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	userResponse, err := s.userService.FindByEmailAndPassword(ctx, req.Email, req.Password)
	if err != nil {
		return nil, apperror.New(apperror.ErrUnauthorized, "Wrong username or password", err)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = strconv.Itoa(userResponse.Id)
	claims["name"] = userResponse.Name
	claims["email"] = userResponse.Email
	claims["role"] = userResponse.Role
	if userResponse.FacultyId != nil {
		claims["facultyId"] = strconv.Itoa(*userResponse.FacultyId)
	}
	claims["exp"] = time.Now().Add(time.Hour * accessTokenTtl).Unix()

	accessToken, err := token.SignedString([]byte(s.config.JwtSecret))
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		AccessToken:  accessToken,
		UserResponse: userResponse,
	}, nil
}
