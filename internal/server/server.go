package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"mcm-api/config"
	"mcm-api/docs"
	"mcm-api/pkg/article"
	"mcm-api/pkg/authz"
	"mcm-api/pkg/comment"
	"mcm-api/pkg/contributesession"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/faculty"
	"mcm-api/pkg/log"
	"mcm-api/pkg/media"
	"mcm-api/pkg/startup"
	"mcm-api/pkg/statistic"
	"mcm-api/pkg/systemdata"
	"mcm-api/pkg/user"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	config            *config.Config
	echo              *echo.Echo
	startupService    *startup.Service
	authHandler       *authz.Handler
	userHandler       *user.Handler
	faculty           *faculty.Handler
	storage           *media.Handler
	contributeSession *contributesession.Handler
	contribution      *contribution.Handler
	article           *article.Handler
	comment           *comment.Handler
	systemdata        *systemdata.Handler
	statistic         *statistic.Handler
}

func newServer(
	config *config.Config,
	startupService *startup.Service,
	authHandler *authz.Handler,
	userHandler *user.Handler,
	facultyHandler *faculty.Handler,
	storage *media.Handler,
	contributeSession *contributesession.Handler,
	contribution *contribution.Handler,
	article *article.Handler,
	comment *comment.Handler,
	systemdata *systemdata.Handler,
	statistic *statistic.Handler,
) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:4200",
			"https://localhost:4200",
			"https://hoppscotch.io",
			config.WebAppUrl,
		},
	}))
	return &Server{
		config:            config,
		echo:              e,
		startupService:    startupService,
		authHandler:       authHandler,
		userHandler:       userHandler,
		faculty:           facultyHandler,
		storage:           storage,
		contributeSession: contributeSession,
		contribution:      contribution,
		article:           article,
		comment:           comment,
		systemdata:        systemdata,
		statistic:         statistic,
	}
}

func (s *Server) registerHandler() {
	s.authHandler.Register(s.echo.Group("auth"))
	s.userHandler.Register(s.echo.Group("users"))
	s.faculty.Register(s.echo.Group("faculties"))
	s.storage.Register(s.echo.Group("storage"))
	s.contributeSession.Register(s.echo.Group("contribute-sessions"))
	s.contribution.Register(s.echo.Group("contributions"))
	s.article.Register(s.echo.Group("articles"))
	s.comment.Register(s.echo.Group("comments"))
	s.systemdata.Register(s.echo.Group("system-data"))
	s.statistic.Register(s.echo.Group("statistics"))
}

// @title 123
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *Server) Swagger() {
	docs.SwaggerInfo.Title = "Magazine Collaboration API"
	docs.SwaggerInfo.Description = "Magazine Collaboration API documentation"
	docs.SwaggerInfo.Version = "1.0"
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Server) Start() {
	err := s.startupService.Run()
	if err != nil {
		log.Logger.Panic("Startup service run failed", zap.Error(err))
	}
	s.Swagger()
	s.registerHandler()
	go func() {
		if err := s.echo.Start(":3000"); err != nil {
			log.Logger.Error("server error", zap.Error(err))
			log.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		log.Logger.Fatal("Error shutting down server", zap.Error(err))
	}
	log.Logger.Info("Bye!")
}
