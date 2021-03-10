package worker

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"mcm-api/config"
	"mcm-api/pkg/article"
	"mcm-api/pkg/converter"
	"mcm-api/pkg/log"
	"mcm-api/pkg/queue"
	"os"
	"os/signal"
	"time"
)

type worker struct {
	cfg            *config.Config
	queue          queue.Queue
	converter      converter.DocumentConverter
	articleService *article.Service
}

func newWorker(
	config *config.Config,
	queue queue.Queue,
	converter converter.DocumentConverter,
	articleService *article.Service,
) *worker {
	return &worker{
		cfg:            config,
		queue:          queue,
		converter:      converter,
		articleService: articleService,
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
			ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Minute*5)
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
	switch message.Topic {
	case queue.ContributionCreated:
		return w.contributionCreatedHandler(ctx, message)
	case queue.ArticleUploaded:
		return w.articleUploadedHandler(ctx, message)
	default:
		return fmt.Errorf("unknown topic %v", message.Topic)
	}
}

func (w worker) contributionCreatedHandler(ctx context.Context, message *queue.Message) error {
	return nil
}

func (w worker) articleUploadedHandler(ctx context.Context, message *queue.Message) error {
	if v, ok := message.Data.(queue.ArticleUploadedPayload); ok {
		result, err := w.converter.Convert(ctx, v.Link, v.User)
		if err != nil {
			return err
		}
		return w.articleService.UpdateLinkPdfForVersion(ctx, v.ArticleId, result.Key)
	} else {
		return errors.New("unknown message")
	}
}
