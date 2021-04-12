package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"go.uber.org/zap"
	"mcm-api/config"
	"mcm-api/pkg/article"
	"mcm-api/pkg/contributesession"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/converter"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/log"
	"mcm-api/pkg/media"
	"mcm-api/pkg/notification"
	"mcm-api/pkg/queue"
	"mcm-api/pkg/user"
	"os"
	"os/signal"
	"time"
)

const JobRuntimeTimeoutMinute = 5

type worker struct {
	cfg                        *config.Config
	queue                      queue.Queue
	converter                  converter.DocumentConverter
	articleService             *article.Service
	notificationService        *notification.Service
	userService                *user.Service
	contributionService        *contribution.Service
	contributionSessionService *contributesession.Service
	mediaService               media.Service
	lock                       *redsync.Redsync
}

func newWorker(
	config *config.Config,
	queue queue.Queue,
	converter converter.DocumentConverter,
	articleService *article.Service,
	notificationService *notification.Service,
	userService *user.Service,
	mediaService media.Service,
	contributionService *contribution.Service,
	contributionSessionService *contributesession.Service,
	lock *redsync.Redsync,
) *worker {
	return &worker{
		cfg:                        config,
		queue:                      queue,
		converter:                  converter,
		articleService:             articleService,
		notificationService:        notificationService,
		userService:                userService,
		contributionService:        contributionService,
		contributionSessionService: contributionSessionService,
		mediaService:               mediaService,
		lock:                       lock,
	}
}

func (w worker) Start() {
	log.Logger.Info("starting worker")
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		s := <-signalChannel
		log.Logger.Info("receive signal", zap.String("signal", s.String()))
		cancelFunc()
		log.Logger.Info("grateful shutdown...")
	}()
poolQueueLoop:
	for {
		select {
		case <-ctx.Done():
			break poolQueueLoop
		default:
			message, err := w.queue.Pop(ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Logger.Error("pop queue error", zap.Error(err))
				}
				break poolQueueLoop
			}

			if message == nil {
				log.Logger.Debug("receive empty message")
				continue
			}
			log.Logger.Info("receive message", zap.Any("message", message))
			ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Minute*JobRuntimeTimeoutMinute)
			err = w.handleMessage(ctxTimeout, message)
			cancelFunc()
			if err != nil {
				log.Logger.Error("process message error", zap.Error(err))
			} else {
				log.Logger.Info("finish process message", zap.Any("message", message))
			}
		}
	}
}

func (w worker) handleMessage(ctx context.Context, message *queue.Message) error {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("recover from panic", zap.Any("error", r))
		}
	}()
	switch message.Topic {
	case queue.ContributionCreated:
		return w.contributionCreatedHandler(ctx, message)
	case queue.ArticleUploaded:
		return w.articleUploadedHandler(ctx, message)
	case queue.ExportContributeSession:
		return w.exportContributeSessionHandler(ctx, message)
	default:
		return fmt.Errorf("unknown topic %v", message.Topic)
	}
}

func (w worker) contributionCreatedHandler(ctx context.Context, message *queue.Message) error {
	if v, ok := message.Data.(*queue.ContributionCreatedPayload); ok {
		entities, err := w.userService.GetAllUserOfFaculty(ctx, enforcer.MarketingCoordinator, v.FacultyId)
		if err != nil {
			return err
		}
		if len(entities) == 0 {
			log.Logger.Info("not found any marketing coordinator to send email")
			return nil
		}
		for _, marketingCoordinator := range entities {
			err = w.notificationService.SendNewContributionEmail(
				&notification.Destination{ToAddresses: []string{marketingCoordinator.Email}},
				&notification.TemplateNewContributionPayLoad{
					Name:        marketingCoordinator.Name,
					StudentName: v.UserName,
					Link:        "https://google.com",
				})
			if err != nil {
				log.Logger.Error("send email failed",
					zap.Error(err),
					zap.Any("target", marketingCoordinator),
				)
			}
		}
		return nil
	} else {
		return errors.New("unknown message")
	}
}

func (w worker) articleUploadedHandler(ctx context.Context, message *queue.Message) error {
	if v, ok := message.Data.(*queue.ArticleUploadedPayload); ok {
		result, err := w.converter.Convert(ctx, v.Link, v.User)
		if err != nil {
			return err
		}
		return w.articleService.UpdateLinkPdfForVersion(ctx, v.ArticleId, result.Key)
	} else {
		return errors.New("unknown message")
	}
}
